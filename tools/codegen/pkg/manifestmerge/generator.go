package manifestmerge

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/utils"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/managedfields"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/spec3"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"k8s.io/utils/pointer"
	"path/filepath"
	kyaml "sigs.k8s.io/yaml"
)

var defaultClusterProfilesToInject = []string{
	"include.release.openshift.io/ibm-cloud-managed",
	"include.release.openshift.io/self-managed-high-availability",
	"include.release.openshift.io/single-node-developer",
}

// Options contains the configuration required for the schemapatch generator.
type Options struct {
	// Disabled indicates whether the schemapatch generator is disabled or not.
	// This default to false as the schemapatch generator is enabled by default.
	Disabled bool

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool

	// TupleOverrides is temporary. Once we generate all CRDs, we won't need this.
	// This allows specification of a CRD that is "ungated". This concept will go away
	// and FeatureGate handling will be required and we'll generate multiple files.
	// It also allows the specification of custom clusterProfiles until we normalize that too.
	TupleOverrides []generation.TupleOverride
}

// generator implements the generation.Generator interface.
// It is designed to generate schemapatch updates for a particular API group.
type generator struct {
	disabled       bool
	verify         bool
	tupleOverrides []generation.TupleOverride
}

// NewGenerator builds a new schemapatch generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		disabled:       opts.Disabled,
		verify:         opts.Verify,
		tupleOverrides: opts.TupleOverrides,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the schemapatch configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.ManifestMerge == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled:       config.ManifestMerge.Disabled,
		Verify:         g.verify,
		TupleOverrides: config.ManifestMerge.TupleOverrides,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "manifestMerge"
}

// GenGroup runs the schemapatch generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	if g.disabled {
		klog.V(2).Infof("Skipping %q for %s", g.Name(), groupCtx.Name)
		return nil
	}

	versionPaths := allVersionPaths(groupCtx.Versions)

	errs := []error{}

	for _, version := range groupCtx.Versions {
		action := "Generating"
		if g.verify {
			action = "Verifying"
		}

		klog.Infof("%s %q for for %s/%s", action, g.Name(), groupCtx.Name, version.Name)

		if err := g.genGroupVersion(groupCtx.Name, version, versionPaths); err != nil {
			errs = append(errs, fmt.Errorf("could not run %q generator for group/version %s/%s: %w", g.Name(), groupCtx.Name, version.Name, err))
		}
	}

	if len(errs) > 0 {
		return kerrors.NewAggregate(errs)
	}

	return nil
}

func (g *generator) findExceptions(crdName string) []generation.TupleOverride {
	ret := []generation.TupleOverride{}
	for _, curr := range g.tupleOverrides {
		if curr.CRDName == crdName {
			ret = append(ret, curr)
		}
	}
	return ret
}

func (g *generator) findExceptionForFeatureSet(crdName, featureSetName string) *generation.TupleOverride {
	for _, curr := range g.tupleOverrides {
		if curr.CRDName == crdName && curr.FeatureSet == featureSetName {
			return &curr
		}
	}
	return nil
}

// genGroupVersion runs the schemapatch generator against a particular version of the API group.
func (g *generator) genGroupVersion(group string, version generation.APIVersionContext, versionPaths []string) error {
	errs := []error{}

	for _, versionPath := range versionPaths {
		resourcePaths := []string{}

		manualCRDOverridesPath := filepath.Join(versionPath, "manual-override-crd-manifests")
		byFeatureGatePath := filepath.Join(versionPath, "zz_generated.featuregated-crd-manifests")
		possibleResources, err := os.ReadDir(byFeatureGatePath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		for _, path := range possibleResources {
			if path.IsDir() {
				resourcePaths = append(resourcePaths, filepath.Join(byFeatureGatePath, path.Name()))
			}
		}

		for _, resourcePath := range resourcePaths {
			// at this point we have a few paths.  In the end, we want to generate manifests for every (crd, clusterprofile, featureset) tuple.
			// a prefix can be specified to associate with an image that creates these, but we create one file each.
			// to start, we need to be able to diff the result against what we already have to do a meaningful review,
			// so to do this we'll allow an all the tuples to a single file option.
			crdName := filepath.Base(resourcePath)

			// allResourcePaths has our generated path and the manual overrides if they exist
			allResourcePaths := []string{resourcePath}
			manualCRDOverridesForCRDPath := filepath.Join(manualCRDOverridesPath, filepath.Base(resourcePath))
			if _, err := os.ReadDir(manualCRDOverridesForCRDPath); err == nil {
				allResourcePaths = append(allResourcePaths, manualCRDOverridesForCRDPath)
			}

			if crdExceptions := g.findExceptions(crdName); len(crdExceptions) == 1 && crdExceptions[0].Ungated {
				// this means we're supposed to create a single file with everything
				resultingCRD, newErrs := mergeAllPertinentCRDsInDirs(allResourcePaths, &AllFeatureGates{}, "")
				if len(newErrs) > 0 {
					errs = append(errs, newErrs...)
					continue
				}
				resultingCRD.SetManagedFields(nil)

				outputFileBaseName := ""
				if outputFilenamePattern := resultingCRD.GetAnnotations()["api.openshift.io/filename-pattern"]; len(outputFilenamePattern) > 0 {
					if !strings.Contains(outputFilenamePattern, "MARKERS") {
						errs = append(errs, fmt.Errorf("crd %q is missing ungated MARKERS from '// +openshift:file-pattern='", crdName))
						continue
					}
					outputFileBaseName = strings.ReplaceAll(outputFilenamePattern, "MARKERS", "")
				}
				if len(outputFileBaseName) == 0 {
					errs = append(errs, fmt.Errorf("crd %q needs '// +openshift:file-pattern=' for ungated", crdName))
					continue
				}
				outputFile := filepath.Join(versionPath, outputFileBaseName)

				annotations := resultingCRD.GetAnnotations()
				delete(annotations, "api.openshift.io/filename-pattern")
				for key := range annotations {
					if strings.HasPrefix(key, "feature-gate.release.openshift.io/") {
						delete(annotations, key)
					}
				}
				if len(crdExceptions[0].ClusterProfilesToInject) > 0 {
					for _, clusterProfile := range crdExceptions[0].ClusterProfilesToInject {
						annotations[clusterProfile] = "true"
					}
				} else {
					for _, clusterProfile := range defaultClusterProfilesToInject {
						annotations[clusterProfile] = "true"
					}
				}
				for key := range annotations {
					if strings.HasSuffix(key, "-") {
						toRemove := key[:len(key)-1]
						delete(annotations, toRemove)
						delete(annotations, key)
					}
				}
				resultingCRD.SetAnnotations(annotations)

				// duplication is ugly, but it's only during the trip from here to there.
				manifestData, err := kyaml.Marshal(resultingCRD)
				if err != nil {
					errs = append(errs, fmt.Errorf("could not encode file %s: %v", outputFile, err))
					continue
				}

				if g.verify {
					existingBytes, err := os.ReadFile(outputFile)
					if err != nil {
						errs = append(errs, fmt.Errorf("could not read file %s: %v", outputFile, err))
						continue
					}
					if !bytes.Equal(manifestData, existingBytes) {
						diff := utils.Diff(existingBytes, manifestData, outputFile)

						return fmt.Errorf("API schema for %s is out of date, please regenerate the API schema:\n%s", outputFile, diff)
					}

					continue
				}

				if err := os.WriteFile(outputFile, manifestData, 0644); err != nil {
					return fmt.Errorf("could not write manifest %s: %w", outputFile, err)
				}

				continue
			}

			// ok, if we're down here, then we need to iterate through all known featuresets
			// again in the future we'll expand to clusterprofile, featureset tuples, but for now all clusterprofiles are considered combined
			// this assumption works for everything *except* for authentication.
			for _, featureSetName := range []string{"Default", "TechPreviewNoUpgrade", "CustomNoUpgrade"} {
				// TODO this will eventually need the clusterprofile too
				partialManifestFilter, err := FilterForFeatureSet(featureSetName)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				var mergeErrors []error
				resultingCRD, mergeErrors := mergeAllPertinentCRDsInDirs(allResourcePaths, partialManifestFilter, featureSetName)
				if len(mergeErrors) > 0 {
					errs = append(errs, mergeErrors...)
					continue
				}

				// TODO the filename is carried on the CRD, need to work out how to clean up. probably easier once we have a dedicated directory
				if resultingCRD == nil { // this means we didn't find any file that matched the filter this is ok, we have nothing to do.
					continue
				}

				outputFileBaseName := ""
				if outputFilenamePattern := resultingCRD.GetAnnotations()["api.openshift.io/filename-pattern"]; len(outputFilenamePattern) > 0 {
					if !strings.Contains(outputFilenamePattern, "MARKERS") {
						errs = append(errs, fmt.Errorf("crd %q is missing featureset/%q MARKERS from '// +openshift:file-pattern=' %q", crdName, featureSetName, outputFilenamePattern))
						continue
					}
					// TODO this should include clusterProfile once we accommodate it.
					fileMarker := fmt.Sprintf("-%s", featureSetName)
					outputFileBaseName = strings.ReplaceAll(outputFilenamePattern, "MARKERS", fileMarker)
				}
				if len(outputFileBaseName) == 0 {
					errs = append(errs, fmt.Errorf("crd %q needs '// +openshift:file-pattern='", crdName))
					continue
				}
				outputFile := filepath.Join(versionPath, outputFileBaseName)

				resultingCRD.SetManagedFields(nil)

				annotations := resultingCRD.GetAnnotations()
				delete(annotations, "api.openshift.io/filename-pattern")
				for key := range annotations {
					if strings.HasPrefix(key, "feature-gate.release.openshift.io/") {
						delete(annotations, key)
					}
				}
				// TODO probably have cases for different clusterprofiles
				exception := g.findExceptionForFeatureSet(crdName, featureSetName)
				if exception != nil && len(exception.ClusterProfilesToInject) > 0 {
					for _, clusterProfile := range exception.ClusterProfilesToInject {
						annotations[clusterProfile] = "true"
					}
				} else {
					for _, clusterProfile := range defaultClusterProfilesToInject {
						annotations[clusterProfile] = "true"
					}
				}
				for key := range annotations {
					if strings.HasSuffix(key, "-") {
						toRemove := key[:len(key)-1]
						delete(annotations, toRemove)
						delete(annotations, key)
					}
				}
				resultingCRD.SetAnnotations(annotations)

				manifestData, err := kyaml.Marshal(resultingCRD.Object)
				if err != nil {
					errs = append(errs, fmt.Errorf("could not encode file %s: %v", outputFile, err))
					continue
				}

				if g.verify {
					existingBytes, err := os.ReadFile(outputFile)
					if err != nil {
						errs = append(errs, fmt.Errorf("could not read file %s: %v", outputFile, err))
						continue
					}
					if !bytes.Equal(manifestData, existingBytes) {
						diff := utils.Diff(existingBytes, manifestData, outputFile)

						return fmt.Errorf("API schema for %s is out of date, please regenerate the API schema:\n%s", outputFile, diff)
					}

					continue
				}

				if err := os.WriteFile(outputFile, manifestData, 0644); err != nil {
					return fmt.Errorf("could not write manifest %s: %w", outputFile, err)
				}
			}

		}
	}

	return kerrors.NewAggregate(errs)
}

// pertinent is determined by the `filter`. If it passes the filter, it's pertinent.
// filters commonly include clusterprofile and featureset-to-feature-gate mapping.
// for example, the TechPreviewNoUpgrade featureset produces a filter that looks to see if the featuregate specified is
// enabled when TechPreviewNoUpgrade is set.
func mergeAllPertinentCRDsInDirs(resourcePaths []string, filter ManifestFilter, featureSet string) (*unstructured.Unstructured, []error) {
	var resultingCRD *unstructured.Unstructured
	var errs []error

	for _, resourcePath := range resourcePaths {
		var currErrs []error
		resultingCRD, currErrs = mergeAllPertinentCRDsInDir(resourcePath, filter, featureSet, resultingCRD)
		errs = append(errs, currErrs...)
	}

	return resultingCRD, errs
}

// pertinent is determined by the `filter`. If it passes the filter, it's pertinent.
// filters commonly include clusterprofile and featureset-to-feature-gate mapping.
// for example, the TechPreviewNoUpgrade featureset produces a filter that looks to see if the featuregate specified is
// enabled when TechPreviewNoUpgrade is set.
func mergeAllPertinentCRDsInDir(resourcePath string, filter ManifestFilter, featureSet string, startingCRD *unstructured.Unstructured) (*unstructured.Unstructured, []error) {

	var unstructuredResultingCRD *unstructured.Unstructured
	if startingCRD != nil {
		unstructuredResultingCRD = startingCRD.DeepCopy()

	} else {
		annotations := map[string]string{
			"api.openshift.io/merged-by-featuregates": "true",
		}
		if len(featureSet) > 0 {
			annotations["release.openshift.io/feature-set"] = featureSet
		}

		unstructuredResultingCRD = &unstructured.Unstructured{}
		unstructuredResultingCRD.GetObjectKind().SetGroupVersionKind(apiextensionsv1.SchemeGroupVersion.WithKind("CustomResourceDefinition"))
		unstructuredResultingCRD.SetAnnotations(annotations)
	}

	errs := []error{}
	var resultingCRD runtime.Object
	resultingCRD = unstructuredResultingCRD

	partialManifestFiles, err := os.ReadDir(resourcePath)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	foundAFile := false
	for _, partialManifest := range partialManifestFiles {
		path := filepath.Join(resourcePath, partialManifest.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file %s: %w", path, err))
			continue
		}
		useManifest, err := filter.UseManifest(data)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not determine whether to use file %s: %w", path, err))
			continue
		}
		if !useManifest {
			continue
		}

		foundAFile = true
		newResult, err := mergeCRD(resultingCRD, data, path)
		if err != nil {
			errs = append(errs, fmt.Errorf("error applying %q: %w", path, err))
			continue
		}
		resultingCRD = newResult
	}

	if !foundAFile {
		return nil, errs
	}

	return resultingCRD.(*unstructured.Unstructured), errs
}

// allVersionPaths creates a list of all version paths for the group.
func allVersionPaths(versions []generation.APIVersionContext) []string {
	out := []string{}

	for _, version := range versions {
		out = append(out, version.Path)
	}

	return out
}

func mergeCRD(obj runtime.Object, patchBytes []byte, fieldManager string) (runtime.Object, error) {
	ssaFieldManager, err := getApplyFieldManager()
	if err != nil {
		return nil, err
	}

	serverSideApplyPatcher := &applyPatcher{
		patch: patchBytes,
		options: &metav1.PatchOptions{
			Force:        pointer.BoolPtr(true),
			FieldManager: fieldManager,
		},
		fieldManager: ssaFieldManager,
	}

	ret := obj.DeepCopyObject()
	output, err := serverSideApplyPatcher.applyPatchToCurrentObject(ret)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// isCustomResourceDefinition returns true if the object is a CustomResourceDefinition.
// This is determined by the object having a Kind of CustomResourceDefinition and the
// correct APIVersion.
func isCustomResourceDefinition(partialObject *metav1.PartialObjectMetadata) bool {
	return partialObject.APIVersion == apiextensionsv1.SchemeGroupVersion.String() && partialObject.Kind == "CustomResourceDefinition"
}

//go:embed crd-schema.json
var crdSchemaJSON []byte

var (
	schemaReadOnce sync.Once
	schemaReadErr  error
	schemaSchema   map[string]*spec.Schema
)

func getCRDSchema() (map[string]*spec.Schema, error) {
	schemaReadOnce.Do(func() {
		openapiV3CRD := &spec3.OpenAPI{}
		err := json.Unmarshal(crdSchemaJSON, openapiV3CRD)
		if err != nil {
			schemaReadErr = err
			return
		}
		schemaSchema = openapiV3CRD.Components.Schemas
	})
	return schemaSchema, schemaReadErr
}

var (
	fieldManagerOnce sync.Once
	fieldManagerErr  error
	fieldManagerRet  *managedfields.FieldManager
)

func getApplyFieldManager() (*managedfields.FieldManager, error) {
	openAPIModels, err := getCRDSchema()
	if err != nil {
		return nil, err
	}

	fieldManagerOnce.Do(func() {
		typeConverter, err := managedfields.NewTypeConverter(openAPIModels, false)
		if err != nil {
			fieldManagerErr = err
			return
		}
		fieldManager, err := managedfields.NewDefaultCRDFieldManager(
			typeConverter,
			noopConverter{},
			noopDefaulter{},
			noopCreator{},
			apiextensionsv1.SchemeGroupVersion.WithKind("CustomResourceDefinition"),
			apiextensionsv1.SchemeGroupVersion,
			"",
			nil,
		)
		if err != nil {
			fieldManagerErr = err
			return
		}
		fieldManagerRet = fieldManager
	})

	return fieldManagerRet, fieldManagerErr
}

type applyPatcher struct {
	patch        []byte
	options      *metav1.PatchOptions
	fieldManager *managedfields.FieldManager
}

func (p *applyPatcher) applyPatchToCurrentObject(obj runtime.Object) (runtime.Object, error) {
	force := false
	if p.options.Force != nil {
		force = *p.options.Force
	}
	if p.fieldManager == nil {
		panic("FieldManager must be installed to run apply")
	}

	patchObj := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := kyaml.Unmarshal(p.patch, &patchObj.Object); err != nil {
		return nil, fmt.Errorf("unable to unmarshal the patch: %w", err)
	}

	obj, err := p.fieldManager.Apply(obj, patchObj, p.options.FieldManager, force)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func readCRDYaml(data []byte) (*unstructured.Unstructured, error) {
	json, err := kyaml.YAMLToJSON(data)
	if err != nil {
		json = data
	}
	obj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, json)
	if err != nil {
		return nil, err
	}

	return obj.(*unstructured.Unstructured), nil
}
