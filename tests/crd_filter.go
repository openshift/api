package tests

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	"path"
	"path/filepath"
	kyaml "sigs.k8s.io/yaml"
	"strings"
)

var (
	clusterProfileToShortName = map[string]string{
		"include.release.openshift.io/hypershift":                     "Hypershift",
		"include.release.openshift.io/ibm-cloud-managed":              "IbmCloudManaged",
		"include.release.openshift.io/self-managed-high-availability": "SelfManagedHA",
		"include.release.openshift.io/single-node-developer":          "SingleNode",
	}
)

func perTestRuntimeInfo(suitePath, crdName, featureGate string) (*PerTestRuntimeInfo, error) {
	requiresFeatureGateDisabled := strings.HasPrefix(featureGate, "-")

	crdFilesToCheck := []string{}

	// account for the generated file move.
	relativePathForCRDs := filepath.Join(suitePath, "..", "..", "zz_generated.crd-manifests")

	generatedCRDs, err := os.ReadDir(relativePathForCRDs)
	if err != nil {
		return nil, err
	}
	for _, currCRDFile := range generatedCRDs {
		relativeFilename := filepath.Join(relativePathForCRDs, currCRDFile.Name())
		filename, err := filepath.Abs(relativeFilename)
		if err != nil {
			return nil, fmt.Errorf("could not generate absolute path for %q: %w", relativeFilename, err)
		}

		currCRD, err := loadCRDFromFile(filename)
		if err != nil {
			// not all files are CRDs, verify will catch garbage.
			continue
		}
		if currCRD.Name != crdName {
			continue
		}
		if len(featureGate) == 0 {
			// test is ungated, check everything
			crdFilesToCheck = append(crdFilesToCheck, filename)
			continue
		}

		featureSet := currCRD.Annotations["release.openshift.io/feature-set"]
		if featureSet == "CustomNoUpgrade" {
			// CustomNoUpgrade includes every field
			if requiresFeatureGateDisabled {
				continue
			}
			crdFilesToCheck = append(crdFilesToCheck, filename)
			continue
		}
		clusterProfilesForCRD := clusterProfilesFrom(currCRD.Annotations)
		if len(clusterProfilesForCRD) == 0 {
			// this is weird, test everything
			crdFilesToCheck = append(crdFilesToCheck, filename)
			continue
		}

		// if the manifest has more than one clusterProfile, then the crd schema must have been the same no matter which
		// featuregates were used.  Simply select the first one to check.
		clusterProfileToCheck := clusterProfilesForCRD.List()[0]
		featureGateStatus, err := featureGatesForClusterProfileFeatureSet("../payload-manifests/featuregates", clusterProfileToCheck, featureSet)
		if err != nil {
			return nil, fmt.Errorf("unable to find featureGates to check for %v: %w", filename, err)
		}

		switch {
		case requiresFeatureGateDisabled:
			disabledFeatureGate := strings.TrimPrefix(featureGate, "-")
			enabled, found := featureGateStatus[disabledFeatureGate]
			if !found {
				return nil, fmt.Errorf("unable to find featureGate/%v to check for %v", featureGate, filename)
			}
			if !enabled {
				crdFilesToCheck = append(crdFilesToCheck, filename)
				continue
			}

		default:
			enabled, found := featureGateStatus[featureGate]
			if !found {
				return nil, fmt.Errorf("unable to find featureGate/%v to check for %v", featureGate, filename)
			}
			if enabled {
				crdFilesToCheck = append(crdFilesToCheck, filename)
			}

		}

	}

	ret := &PerTestRuntimeInfo{
		CRDFilenames: crdFilesToCheck,
	}
	return ret, nil
}

func clusterProfilesFrom(annotations map[string]string) sets.String {
	ret := sets.NewString()
	for k, v := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") && v == "true" {
			ret.Insert(k)
		}
	}
	return ret
}

func clusterProfilesShortNamesFrom(annotations map[string]string) sets.String {
	ret := sets.NewString()
	for k, v := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") && v == "true" {
			ret.Insert(clusterProfileToShortName[k])
		}
	}
	return ret
}

func featureGatesForClusterProfileFeatureSet(payloadFeatureGatePath, clusterProfile, featureSetName string) (map[string]bool, error) {
	if len(featureSetName) == 0 {
		// if the featureSetName is ungated, then all CRD schemas for every featureset on this clusterProfile must be the same.
		// Choose Default so that we get a valid manifest to check.
		featureSetName = "Default"
	}

	featureGateFilename := path.Join(
		payloadFeatureGatePath,
		fmt.Sprintf("featureGate-%s-%s.yaml", clusterProfileToShortName[clusterProfile], featureSetName),
	)
	featureGateBytes, err := os.ReadFile(featureGateFilename)
	if err != nil {
		return nil, err
	}
	// use unstructured to pull this information to avoid vendoring openshift/api
	uncastFeatureGate := map[string]interface{}{}
	if err := kyaml.Unmarshal(featureGateBytes, &uncastFeatureGate); err != nil {
		return nil, fmt.Errorf("unable to parse featuregate %q: %w", featureGateFilename, err)
	}

	uncastFeatureGateSlice, _, err := unstructured.NestedSlice(uncastFeatureGate, "status", "featureGates")
	if err != nil {
		return nil, fmt.Errorf("no slice found %w", err)
	}
	enabledFeatureGates, _, err := unstructured.NestedSlice(uncastFeatureGateSlice[0].(map[string]interface{}), "enabled")
	if err != nil {
		return nil, fmt.Errorf("no enabled found %w", err)
	}
	disabledFeatureGates, _, err := unstructured.NestedSlice(uncastFeatureGateSlice[0].(map[string]interface{}), "disabled")
	if err != nil {
		return nil, fmt.Errorf("no enabled found %w", err)
	}

	featureGateMapping := map[string]bool{}
	for _, currGate := range enabledFeatureGates {
		featureGateName, _, err := unstructured.NestedString(currGate.(map[string]interface{}), "name")
		if err != nil {
			return nil, fmt.Errorf("no gate name found %w", err)
		}
		featureGateMapping[featureGateName] = true
	}
	for _, currGate := range disabledFeatureGates {
		featureGateName, _, err := unstructured.NestedString(currGate.(map[string]interface{}), "name")
		if err != nil {
			return nil, fmt.Errorf("no gate name found %w", err)
		}
		featureGateMapping[featureGateName] = false
	}

	return featureGateMapping, nil
}
