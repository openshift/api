package schemapatch

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	yamlpatch "github.com/vmware-archive/yaml-patch"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

// executeYAMLPatchForGroupVersion runs the yaml patch against a particular version.
func executeYAMLPatchForGroupVersion(version generation.APIVersionContext, requiredFeatureSets []sets.String) error {
	patchPairs, err := loadPatchFilesForGroupVersion(version.Path)
	if err != nil {
		return fmt.Errorf("could not load patch files: %v", err)
	}

	for _, pair := range patchPairs {
		if err := patchFile(version.Path, pair.baseFile, pair.patchFile, requiredFeatureSets); err != nil {
			return fmt.Errorf("could not patch file: %v", err)
		}
	}

	return nil
}

// patchPair pairs the patch file to it's base CRD.
type patchPair struct {
	patchFile fs.FileInfo
	baseFile  fs.FileInfo
}

// loadPatchFilesForGroupVersion looks for patch pairs within the given group version directory.
func loadPatchFilesForGroupVersion(groupVersionDir string) ([]patchPair, error) {
	out := []patchPair{}

	dirEntries, err := ioutil.ReadDir(groupVersionDir)
	if err != nil {
		return nil, fmt.Errorf("could not read file info for directory %s: %v", groupVersionDir, err)
	}

	for _, fileInfo := range dirEntries {
		// Find all files that are yaml-patches
		if fileInfo.IsDir() || filepath.Ext(fileInfo.Name()) != ".yaml-patch" {
			continue
		}

		baseCRDName := strings.TrimRight(fileInfo.Name(), "-patch")
		baseCRDInfo, err := findFile(dirEntries, baseCRDName)
		if err != nil {
			return nil, fmt.Errorf("could not find base CRD file for patch %s: %v", fileInfo.Name(), err)
		}

		out = append(out, patchPair{
			patchFile: fileInfo,
			baseFile:  baseCRDInfo,
		})
	}

	return out, nil
}

// findFile looks for a file with the given name within the directory entries passed.
func findFile(dirEntries []fs.FileInfo, name string) (fs.FileInfo, error) {
	for _, fileInfo := range dirEntries {
		if fileInfo.Name() == name {
			return fileInfo, nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", name)
}

// patchFile applies a patch file to a CRD file.
func patchFile(groupVersionDir string, baseInfo, patchInfo fs.FileInfo, requiredFeatureSets []sets.String) error {
	klog.V(2).Infof("Patching CRD %s with patch file %s", baseInfo.Name(), patchInfo.Name())

	placeholderWrapper := yamlpatch.NewPlaceholderWrapper("{{", "}}")

	patchPath := filepath.Join(groupVersionDir, patchInfo.Name())
	patch, err := loadPatch(placeholderWrapper, patchPath)
	if err != nil {
		return fmt.Errorf("could not load patch %s: %v", patchPath, err)
	}

	basePath := filepath.Join(groupVersionDir, baseInfo.Name())
	baseDoc, err := loadBaseDoc(placeholderWrapper, basePath)
	if err != nil {
		return fmt.Errorf("could not load base %s: %v", basePath, err)
	}

	if !mayHandleFile(baseDoc, requiredFeatureSets) {
		klog.V(3).Infof("Skipping patch %s as it is not required for the current feature sets", baseInfo.Name())
		return nil
	}

	patchedDoc, err := patch.Apply(baseDoc)
	if err != nil {
		return fmt.Errorf("could not apply patch: %v", err)
	}

	if err := writePatchedFile(placeholderWrapper, basePath, baseInfo.Mode(), patchedDoc); err != nil {
		return fmt.Errorf("could not write file %s: %v", basePath, err)
	}

	return nil
}

// loadPatch loads and parses the patch file from disk.
func loadPatch(placeholderWrapper *yamlpatch.PlaceholderWrapper, path string) (yamlpatch.Patch, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	patch, err := yamlpatch.DecodePatch(placeholderWrapper.Wrap(data))
	if err != nil {
		return nil, fmt.Errorf("could not decode patch: %v", err)
	}

	return patch, nil
}

// loadBaseDoc loads the base CRD document from disk.
func loadBaseDoc(placeholderWrapper *yamlpatch.PlaceholderWrapper, path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	return placeholderWrapper.Wrap(data), nil
}

// writePatchedFile writes the patched CRD document to disk.
func writePatchedFile(placeholderWrapper *yamlpatch.PlaceholderWrapper, name string, mode fs.FileMode, data []byte) error {
	data = placeholderWrapper.Unwrap(data)

	if err := ioutil.WriteFile(name, data, mode); err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	return nil
}
