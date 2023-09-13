package manifestcomparators

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type noFieldRemoval struct{}

func NoFieldRemoval() CRDComparator {
	return noFieldRemoval{}
}

func (noFieldRemoval) Name() string {
	return "NoFieldRemoval"
}

func (noFieldRemoval) WhyItMatters() string {
	return "If fields are removed, then clients that rely on those fields will not be able to read them or write them."
}

func (b noFieldRemoval) Compare(existingCRD, newCRD *apiextensionsv1.CustomResourceDefinition) (ComparisonResults, error) {
	if existingCRD == nil {
		return ComparisonResults{
			Name:         b.Name(),
			WhyItMatters: b.WhyItMatters(),

			Errors:   nil,
			Warnings: nil,
			Infos:    nil,
		}, nil
	}
	errsToReport := []string{}

	for _, newVersion := range newCRD.Spec.Versions {

		existingVersion := GetVersionByName(existingCRD, newVersion.Name)
		if existingVersion == nil {
			continue
		}

		existingFields := sets.NewString()
		SchemaHas(existingVersion.Schema.OpenAPIV3Schema, field.NewPath("^"), field.NewPath("^"), func(s *apiextensionsv1.JSONSchemaProps, fldPath, simpleLocation *field.Path) bool {
			existingFields.Insert(simpleLocation.String())
			return false
		})

		newFields := sets.NewString()
		SchemaHas(newVersion.Schema.OpenAPIV3Schema, field.NewPath("^"), field.NewPath("^"), func(s *apiextensionsv1.JSONSchemaProps, fldPath, simpleLocation *field.Path) bool {
			newFields.Insert(simpleLocation.String())
			return false
		})

		removedFields := existingFields.Difference(newFields)
		for _, removedField := range removedFields.List() {
			errsToReport = append(errsToReport, fmt.Sprintf("crd/%v version/%v field/%v may not be removed", newCRD.Name, newVersion.Name, removedField))
		}

	}

	return ComparisonResults{
		Name:         b.Name(),
		WhyItMatters: b.WhyItMatters(),

		Errors:   errsToReport,
		Warnings: nil,
		Infos:    nil,
	}, nil
}
