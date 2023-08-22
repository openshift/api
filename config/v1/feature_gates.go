package v1

// FeatureGateDescription is a golang-only interface used to contains details for a feature gate.
type FeatureGateDescription struct {
	// FeatureGateAttributes is the information that appears in the API
	FeatureGateAttributes FeatureGateAttributes

	// OwningJiraComponent is the jira component that owns most of the impl and first assignment for the bug.
	// This is the team that owns the feature long term.
	OwningJiraComponent string
	// ResponsiblePerson is the person who is on the hook for first contact.  This is often, but not always, a team lead.
	// It is someone who can make the promise on the behalf of the team.
	ResponsiblePerson string
	// OwningProduct is the product that owns the lifecycle of the gate.
	OwningProduct OwningProduct
}

type OwningProduct string

var (
	ocpSpecific = OwningProduct("OCP")
	kubernetes  = OwningProduct("Kubernetes")
)

var (
	FeatureGateValidatingAdmissionPolicy = FeatureGateName("ValidatingAdmissionPolicy")
	ValidatingAdmissionPolicy            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateValidatingAdmissionPolicy,
		},
		OwningJiraComponent: "kube-apiserver",
		ResponsiblePerson:   "benluddy",
		OwningProduct:       kubernetes,
	}

	FeatureGateGatewayAPI = FeatureGateName("GatewayAPI")
	GateGatewayAPI        = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateGatewayAPI,
		},
		OwningJiraComponent: "Routing",
		ResponsiblePerson:   "miciah",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateOpenShiftPodSecurityAdmission = FeatureGateName("OpenShiftPodSecurityAdmission")
	OpenShiftPodSecurityAdmission            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateOpenShiftPodSecurityAdmission,
		},
		OwningJiraComponent: "auth",
		ResponsiblePerson:   "stlaz",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateRetroactiveDefaultStorageClass = FeatureGateName("RetroactiveDefaultStorageClass")
	RetroactiveDefaultStorageClass            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateRetroactiveDefaultStorageClass,
		},
		OwningJiraComponent: "storage",
		ResponsiblePerson:   "RomanBednar",
		OwningProduct:       kubernetes,
	}

	FeatureGateExternalCloudProvider = FeatureGateName("ExternalCloudProvider")
	ExternalCloudProvider            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateExternalCloudProvider,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "jspeed",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateExternalCloudProviderAzure = FeatureGateName("ExternalCloudProviderAzure")
	ExternalCloudProviderAzure            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateExternalCloudProviderAzure,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "jspeed",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateExternalCloudProviderGCP = FeatureGateName("ExternalCloudProviderGCP")
	ExternalCloudProviderGCP            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateExternalCloudProviderGCP,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "jspeed",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateExternalCloudProviderExternal = FeatureGateName("ExternalCloudProviderExternal")
	ExternalCloudProviderExternal            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateExternalCloudProviderExternal,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "elmiko",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateCSIDriverSharedResource = FeatureGateName("CSIDriverSharedResource")
	CSIDriverSharedResource            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateCSIDriverSharedResource,
		},
		OwningJiraComponent: "builds",
		ResponsiblePerson:   "adkaplan",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateBuildCSIVolumes = FeatureGateName("BuildCSIVolumes")
	BuildCSIVolumes            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateBuildCSIVolumes,
		},
		OwningJiraComponent: "builds",
		ResponsiblePerson:   "adkaplan",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateNodeSwap = FeatureGateName("NodeSwap")
	NodeSwap            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateNodeSwap,
		},
		OwningJiraComponent: "node",
		ResponsiblePerson:   "ehashman",
		OwningProduct:       kubernetes,
	}

	FeatureGateMachineAPIProviderOpenStack = FeatureGateName("MachineAPIProviderOpenStack")
	MachineAPIProviderOpenStack            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateMachineAPIProviderOpenStack,
		},
		OwningJiraComponent: "openstack",
		ResponsiblePerson:   "egarcia",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateInsightsConfigAPI = FeatureGateName("InsightsConfigAPI")
	InsightsConfigAPI            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateInsightsConfigAPI,
		},
		OwningJiraComponent: "insights",
		ResponsiblePerson:   "tremes",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateDynamicResourceAllocation = FeatureGateName("DynamicResourceAllocation")
	DynamicResourceAllocation            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateDynamicResourceAllocation,
		},
		OwningJiraComponent: "scheduling",
		ResponsiblePerson:   "jchaloup",
		OwningProduct:       kubernetes,
	}

	FeatureGateAdmissionWebhookMatchConditions = FeatureGateName("AdmissionWebhookMatchConditions")
	AdmissionWebhookMatchConditions            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateAdmissionWebhookMatchConditions,
		},
		OwningJiraComponent: "kube-apiserver",
		ResponsiblePerson:   "benluddy",
		OwningProduct:       kubernetes,
	}

	FeatureGateAzureWorkloadIdentity = FeatureGateName("AzureWorkloadIdentity")
	AzureWorkloadIdentity            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateAzureWorkloadIdentity,
		},
		OwningJiraComponent: "cloud-credential-operator",
		ResponsiblePerson:   "abutcher",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateMaxUnavailableStatefulSet = FeatureGateName("MaxUnavailableStatefulSet")
	MaxUnavailableStatefulSet            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateMaxUnavailableStatefulSet,
		},
		OwningJiraComponent: "apps",
		ResponsiblePerson:   "atiratree",
		OwningProduct:       kubernetes,
	}

	FeatureGateEventedPLEG = FeatureGateName("EventedPLEG")
	EventedPleg            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateEventedPLEG,
		},
		OwningJiraComponent: "node",
		ResponsiblePerson:   "sairameshv",
		OwningProduct:       kubernetes,
	}

	FeatureGatePrivateHostedZoneAWS = FeatureGateName("PrivateHostedZoneAWS")
	PrivateHostedZoneAWS            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGatePrivateHostedZoneAWS,
		},
		OwningJiraComponent: "Routing",
		ResponsiblePerson:   "miciah",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateSigstoreImageVerification = FeatureGateName("SigstoreImageVerification")
	SigstoreImageVerification            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateSigstoreImageVerification,
		},
		OwningJiraComponent: "node",
		ResponsiblePerson:   "sgrunert",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateGCPLabelsTags = FeatureGateName("GCPLabelsTags")
	GCPLabelsTags            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateGCPLabelsTags,
		},
		OwningJiraComponent: "Installer",
		ResponsiblePerson:   "bhb",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateAlibabaPlatform = FeatureGateName("AlibabaPlatform")
	AlibabaPlatform            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateAlibabaPlatform,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "jspeed",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateCloudDualStackNodeIPs = FeatureGateName("CloudDualStackNodeIPs")
	CloudDualStackNodeIPs            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateCloudDualStackNodeIPs,
		},
		OwningJiraComponent: "machine-config-operator/platform-baremetal",
		ResponsiblePerson:   "mkowalsk",
		OwningProduct:       kubernetes,
	}
	FeatureGateVSphereStaticIPs = FeatureGateName("VSphereStaticIPs")
	VSphereStaticIPs            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateVSphereStaticIPs,
		},
		OwningJiraComponent: "splat",
		ResponsiblePerson:   "rvanderp3",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateRouteExternalCertificate = FeatureGateName("RouteExternalCertificate")
	RouteExternalCertificate            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateRouteExternalCertificate,
		},
		OwningJiraComponent: "router",
		ResponsiblePerson:   "thejasn",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateAdminNetworkPolicy = FeatureGateName("AdminNetworkPolicy")
	AdminNetworkPolicy            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateAdminNetworkPolicy,
		},
		OwningJiraComponent: "Networking/ovn-kubernetes",
		ResponsiblePerson:   "tssurya",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateAutomatedEtcdBackup = FeatureGateName("AutomatedEtcdBackup")
	AutomatedEtcdBackup            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateAutomatedEtcdBackup,
		},
		OwningJiraComponent: "etcd",
		ResponsiblePerson:   "hasbro17",
		OwningProduct:       ocpSpecific,
	}

	FeatureGateMachineAPIOperatorDisableMachineHealthCheckController = FeatureGateName("MachineAPIOperatorDisableMachineHealthCheckController")
	MachineAPIOperatorDisableMachineHealthCheckController            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: FeatureGateMachineAPIOperatorDisableMachineHealthCheckController,
		},
		OwningJiraComponent: "ecoproject",
		ResponsiblePerson:   "msluiter",
		OwningProduct:       ocpSpecific,
	}
)
