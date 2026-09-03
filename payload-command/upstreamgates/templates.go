package upstreamgates

const FeatureGateBuilderTemplate =`
	FeatureGate%s = newFeatureGate("%s").
						reportProblemsToJiraComponent("%s").
						contactPerson("%s").
						productScope(kubernetes).
						enhancementPR(%s).
						enable(%s).
						mustRegister()
`

const FeatureGateFileTemplate = `package features

var (
%s
)
`
