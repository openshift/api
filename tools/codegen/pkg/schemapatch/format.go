package schemapatch

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/openshift/api/tools/codegen/pkg/generation"

	yaml "gopkg.in/yaml.v3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "sigs.k8s.io/yaml"
)

// formatManifestsForGroup looks for CRD manifests within the group version and formats them
// with a yaml formatter.
func formatManifestsForGroupVersion(version generation.APIVersionContext, requiredFeatureSets sets.String) error {
	errs := []error{}

	dirEntries, err := ioutil.ReadDir(version.Path)
	if err != nil {
		return fmt.Errorf("could not read file info for directory %s: %v", version.Path, err)
	}

	for _, fileInfo := range dirEntries {
		// Find all files that are yaml-patches
		if fileInfo.IsDir() || filepath.Ext(fileInfo.Name()) != ".yaml" {
			continue
		}

		filePath := filepath.Join(version.Path, fileInfo.Name())
		if err := formatManifest(filePath, fileInfo.Mode(), requiredFeatureSets); err != nil {
			errs = append(errs, fmt.Errorf("could not format file %s: %v", fileInfo.Name(), err))
		}
	}

	if len(errs) > 0 {
		return kerrors.NewAggregate(errs)
	}

	return nil
}

// formatManifest formats a particular file. It checks that the file is a CRD before
// formatting it.
// Indentation is pinned to 2 as this is what other kube like tooling uses.
func formatManifest(path string, mode fs.FileMode, requiredFeatureSets sets.String) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", path, err)
	}

	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return fmt.Errorf("could not unmarshal YAML for type meta inspection: %v", err)
	}

	if partialObject.APIVersion != apiextensionsv1.SchemeGroupVersion.String() || partialObject.Kind != "CustomResourceDefinition" || !mayHandleObject(partialObject, requiredFeatureSets) {
		return nil
	}

	node := &yaml.Node{}
	if err := yaml.Unmarshal(data, node); err != nil {
		return fmt.Errorf("could not unmarshal YAML: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)

	enc.SetIndent(2)

	if err := enc.Encode(node); err != nil {
		return fmt.Errorf("could not encode YAML: %v", err)
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), mode); err != nil {
		return fmt.Errorf("could not write file %s: %v", path, err)
	}

	return nil
}
