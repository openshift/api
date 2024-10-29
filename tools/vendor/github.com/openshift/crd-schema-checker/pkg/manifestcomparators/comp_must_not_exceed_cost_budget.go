package manifestcomparators

import (
	"fmt"
	"math"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsvalidation "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema/cel"
	"k8s.io/apimachinery/pkg/util/validation/field"
	celconfig "k8s.io/apiserver/pkg/apis/cel"
	apiservercel "k8s.io/apiserver/pkg/cel"
	"k8s.io/apiserver/pkg/cel/environment"
)

type mustNotExceedCostBudget struct{}

func MustNotExceedCostBudget() CRDComparator {
	return mustNotExceedCostBudget{}
}

func (mustNotExceedCostBudget) Name() string {
	return "MustNotExceedCostBudget"
}

func (mustNotExceedCostBudget) WhyItMatters() string {
	return ""
}

func (b mustNotExceedCostBudget) Validate(crd *apiextensionsv1.CustomResourceDefinition) (ComparisonResults, error) {
	errsToReport := []string{}

	for _, newVersion := range crd.Spec.Versions {
		schema := &apiextensions.JSONSchemaProps{}
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(newVersion.Schema.OpenAPIV3Schema, schema, nil); err != nil {
			errsToReport = append(errsToReport, err.Error())
			continue
		}

		rootCELContext := apiextensionsvalidation.RootCELContext(schema)

		SchemaHas(newVersion.Schema.OpenAPIV3Schema, field.NewPath("^"), field.NewPath("^"), nil,
			func(s *apiextensionsv1.JSONSchemaProps, fldPath, simpleLocation *field.Path, ancestry []*apiextensionsv1.JSONSchemaProps) bool {
				schema := &apiextensions.JSONSchemaProps{}
				if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(s, schema, nil); err != nil {
					errsToReport = append(errsToReport, err.Error())
					return false
				}

				celContext, err := extractCELContext(append(ancestry, s))
				if err != nil {
					errsToReport = append(errsToReport, err.Error())
					return false
				}

				typeInfo, err := celContext.TypeInfo()
				if err != nil {
					errsToReport = append(errsToReport, err.Error())
					return false
				}

				if typeInfo == nil {
					// No validations to check.
					return false
				}

				compResults, err := cel.Compile(
					typeInfo.Schema,
					typeInfo.DeclType,
					celconfig.PerCallLimit,
					environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion(), true),
					cel.NewExpressionsEnvLoader(),
				)
				if err != nil {
					fieldErr := field.InternalError(fldPath, fmt.Errorf("failed to compile x-kubernetes-validations rules: %w", err))
					errsToReport = append(errsToReport, fieldErr.Error())
					return false
				}

				for i, cr := range compResults {
					expressionCost := getExpressionCost(cr, celContext)
					if expressionCost > apiextensionsvalidation.StaticEstimatedCostLimit {
						costErrorMsg := getCostErrorMessage("estimated rule cost", expressionCost, apiextensionsvalidation.StaticEstimatedCostLimit)
						errsToReport = append(errsToReport, field.Forbidden(fldPath, costErrorMsg).Error())
					}
					if rootCELContext.TotalCost != nil {
						rootCELContext.TotalCost.ObserveExpressionCost(fldPath, expressionCost)
					}
					if cr.Error != nil {
						if cr.Error.Type == apiservercel.ErrorTypeRequired {
							errsToReport = append(errsToReport, field.Required(fldPath, cr.Error.Detail).Error())
						} else {
							errsToReport = append(errsToReport, field.Invalid(fldPath, schema.XValidations[i], cr.Error.Detail).Error())
						}
					}
					if cr.MessageExpressionError != nil {
						errsToReport = append(errsToReport, field.Invalid(fldPath, schema.XValidations[i], cr.MessageExpressionError.Detail).Error())
					} else if cr.MessageExpression != nil {
						if cr.MessageExpressionMaxCost > apiextensionsvalidation.StaticEstimatedCostLimit {
							costErrorMsg := getCostErrorMessage("estimated messageExpression cost", cr.MessageExpressionMaxCost, apiextensionsvalidation.StaticEstimatedCostLimit)
							errsToReport = append(errsToReport, field.Forbidden(fldPath, costErrorMsg).Error())
						}
						if celContext.TotalCost != nil {
							celContext.TotalCost.ObserveExpressionCost(fldPath, cr.MessageExpressionMaxCost)
						}
					}
				}

				return false
			})
	}

	return ComparisonResults{
		Name:         b.Name(),
		WhyItMatters: b.WhyItMatters(),

		Errors:   errsToReport,
		Warnings: nil,
		Infos:    nil,
	}, nil
}

func (b mustNotExceedCostBudget) Compare(existingCRD, newCRD *apiextensionsv1.CustomResourceDefinition) (ComparisonResults, error) {
	return RatchetCompare(b, existingCRD, newCRD)
}

// multiplyWithOverflowGuard returns the product of baseCost and cardinality unless that product
// would exceed math.MaxUint, in which case math.MaxUint is returned.
func multiplyWithOverflowGuard(baseCost, cardinality uint64) uint64 {
	if baseCost == 0 {
		// an empty rule can return 0, so guard for that here
		return 0
	} else if math.MaxUint/baseCost < cardinality {
		return math.MaxUint
	}
	return baseCost * cardinality
}

// unbounded uses nil to represent an unbounded cardinality value.
var unbounded *uint64 = nil //nolint:revive // Using as a named variable to provide the meaning of nil in this context.

func getExpressionCost(cr cel.CompilationResult, cardinalityCost *apiextensionsvalidation.CELSchemaContext) uint64 {
	if cardinalityCost.MaxCardinality != unbounded {
		return multiplyWithOverflowGuard(cr.MaxCost, *cardinalityCost.MaxCardinality)
	}
	return multiplyWithOverflowGuard(cr.MaxCost, cr.MaxCardinality)
}

func getCostErrorMessage(costName string, expressionCost, costLimit uint64) string {
	exceedFactor := float64(expressionCost) / float64(costLimit)
	var factor string
	if exceedFactor > 100.0 {
		// if exceedFactor is greater than 2 orders of magnitude, the rule is likely O(n^2) or worse
		// and will probably never validate without some set limits
		// also in such cases the cost estimation is generally large enough to not add any value
		factor = "more than 100x"
	} else if exceedFactor < 1.5 {
		factor = fmt.Sprintf("%fx", exceedFactor) // avoid reporting "exceeds budge by a factor of 1.0x"
	} else {
		factor = fmt.Sprintf("%.1fx", exceedFactor)
	}
	return fmt.Sprintf("%s exceeds budget by factor of %s (try simplifying the rule, or adding maxItems, maxProperties, and maxLength where arrays, maps, and strings are declared)", costName, factor)
}

// extractCELContext takes a series of CEL contextxs and returns the child context of the last schema in the series.
// This is used so that the calculated maximum cardinality of the field is correct.
func extractCELContext(schemas []*apiextensionsv1.JSONSchemaProps) (*apiextensionsvalidation.CELSchemaContext, error) {
	var celContext *apiextensionsvalidation.CELSchemaContext

	for _, s := range schemas {
		schema := &apiextensions.JSONSchemaProps{}
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(s, schema, nil); err != nil {
			return nil, fmt.Errorf("failed to convert schema: %w", err)
		}

		if celContext == nil {
			celContext = apiextensionsvalidation.RootCELContext(schema)
			continue
		}

		celContext = celContext.ChildPropertyContext(schema, s.ID)
	}

	return celContext, nil
}
