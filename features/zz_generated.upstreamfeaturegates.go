package features

var (

	FeatureGateAPIResponseCompression = newFeatureGate("APIResponseCompression").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAPIServerIdentity = newFeatureGate("APIServerIdentity").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAPIServingWithRoutine = newFeatureGate("APIServingWithRoutine").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAllowParsingUserUIDFromCertAuth = newFeatureGate("AllowParsingUserUIDFromCertAuth").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAllowServiceExternalIPs = newFeatureGate("AllowServiceExternalIPs").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAllowUnsafeMalformedObjectDeletion = newFeatureGate("AllowUnsafeMalformedObjectDeletion").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAtomicFIFO = newFeatureGate("AtomicFIFO").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateAuthorizePodWebsocketUpgradeCreatePermission = newFeatureGate("AuthorizePodWebsocketUpgradeCreatePermission").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCBORServingAndStorage = newFeatureGate("CBORServingAndStorage").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCPUManagerPolicyAlphaOptions = newFeatureGate("CPUManagerPolicyAlphaOptions").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCPUManagerPolicyBetaOptions = newFeatureGate("CPUManagerPolicyBetaOptions").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCRDObservedGenerationTracking = newFeatureGate("CRDObservedGenerationTracking").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCRIListStreaming = newFeatureGate("CRIListStreaming").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCSIVolumeHealth = newFeatureGate("CSIVolumeHealth").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClearingNominatedNodeNameAfterBinding = newFeatureGate("ClearingNominatedNodeNameAfterBinding").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClientsAllowCARotation = newFeatureGate("ClientsAllowCARotation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClientsAllowCBOR = newFeatureGate("ClientsAllowCBOR").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClientsAllowTLSCacheGC = newFeatureGate("ClientsAllowTLSCacheGC").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClientsPreferCBOR = newFeatureGate("ClientsPreferCBOR").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCloudControllerManagerWatchBasedRoutesReconciliation = newFeatureGate("CloudControllerManagerWatchBasedRoutesReconciliation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCloudControllerManagerWebhook = newFeatureGate("CloudControllerManagerWebhook").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClusterTrustBundle = newFeatureGate("ClusterTrustBundle").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateClusterTrustBundleProjection = newFeatureGate("ClusterTrustBundleProjection").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateComponentFlagz = newFeatureGate("ComponentFlagz").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateComponentStatusz = newFeatureGate("ComponentStatusz").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateConcurrentWatchObjectDecode = newFeatureGate("ConcurrentWatchObjectDecode").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateConstrainedImpersonation = newFeatureGate("ConstrainedImpersonation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateContainerCheckpoint = newFeatureGate("ContainerCheckpoint").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateContainerRestartRules = newFeatureGate("ContainerRestartRules").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateContainerStopSignals = newFeatureGate("ContainerStopSignals").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCoordinatedLeaderElection = newFeatureGate("CoordinatedLeaderElection").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCrossNamespaceVolumeDataSource = newFeatureGate("CrossNamespaceVolumeDataSource").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateCustomCPUCFSQuotaPeriod = newFeatureGate("CustomCPUCFSQuotaPeriod").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAConsumableCapacity = newFeatureGate("DRAConsumableCapacity").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRADeviceBindingConditions = newFeatureGate("DRADeviceBindingConditions").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRADeviceTaintRules = newFeatureGate("DRADeviceTaintRules").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade(),withGroupKindResources(groupKindResource{Group: "resource.k8s.io", Kind: "DeviceTaintRule", Resource: "devicetaintrules"})).
						mustRegister()


	FeatureGateDRADeviceTaints = newFeatureGate("DRADeviceTaints").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAExtendedResource = newFeatureGate("DRAExtendedResource").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAListTypeAttributes = newFeatureGate("DRAListTypeAttributes").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRANodeAllocatableResources = newFeatureGate("DRANodeAllocatableResources").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAPartitionableDevices = newFeatureGate("DRAPartitionableDevices").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAPrioritizedList = newFeatureGate("DRAPrioritizedList").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAResourceClaimDeviceStatus = newFeatureGate("DRAResourceClaimDeviceStatus").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAResourceClaimGranularStatusAuthorization = newFeatureGate("DRAResourceClaimGranularStatusAuthorization").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAResourcePoolStatus = newFeatureGate("DRAResourcePoolStatus").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRASchedulerFilterTimeout = newFeatureGate("DRASchedulerFilterTimeout").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDRAWorkloadResourceClaims = newFeatureGate("DRAWorkloadResourceClaims").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDeclarativeValidationBeta = newFeatureGate("DeclarativeValidationBeta").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDeploymentReplicaSetTerminatingReplicas = newFeatureGate("DeploymentReplicaSetTerminatingReplicas").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateDetectCacheInconsistency = newFeatureGate("DetectCacheInconsistency").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateEnvFiles = newFeatureGate("EnvFiles").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateEventedPLEG = newFeatureGate("EventedPLEG").
						reportProblemsToJiraComponent("node").
						contactPerson("sairameshv").
						productScope(kubernetes).
						enhancementPR("https://github.com/kubernetes/enhancements/issues/3386").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateExtendWebSocketsToKubelet = newFeatureGate("ExtendWebSocketsToKubelet").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateGangScheduling = newFeatureGate("GangScheduling").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateGenericWorkload = newFeatureGate("GenericWorkload").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateGracefulNodeShutdown = newFeatureGate("GracefulNodeShutdown").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateGracefulNodeShutdownBasedOnPodPriority = newFeatureGate("GracefulNodeShutdownBasedOnPodPriority").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateHPAConfigurableTolerance = newFeatureGate("HPAConfigurableTolerance").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateHPAScaleToZero = newFeatureGate("HPAScaleToZero").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateHostnameOverride = newFeatureGate("HostnameOverride").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateImageVolumeWithDigest = newFeatureGate("ImageVolumeWithDigest").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInOrderInformersBatchProcess = newFeatureGate("InOrderInformersBatchProcess").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInPlacePodLevelResourcesVerticalScaling = newFeatureGate("InPlacePodLevelResourcesVerticalScaling").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInPlacePodVerticalScalingExclusiveCPUs = newFeatureGate("InPlacePodVerticalScalingExclusiveCPUs").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInPlacePodVerticalScalingExclusiveMemory = newFeatureGate("InPlacePodVerticalScalingExclusiveMemory").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInPlacePodVerticalScalingInitContainers = newFeatureGate("InPlacePodVerticalScalingInitContainers").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateInformerResourceVersion = newFeatureGate("InformerResourceVersion").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKMSv1 = newFeatureGate("KMSv1").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("dgrisonnet").
						productScope(kubernetes).
						enhancementPR(legacyFeatureGateWithoutEnhancement).
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKubeletCrashLoopBackOffMax = newFeatureGate("KubeletCrashLoopBackOffMax").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKubeletEnsureSecretPulledImages = newFeatureGate("KubeletEnsureSecretPulledImages").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKubeletInUserNamespace = newFeatureGate("KubeletInUserNamespace").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKubeletSeparateDiskGC = newFeatureGate("KubeletSeparateDiskGC").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateKubeletServiceAccountTokenForCredentialProviders = newFeatureGate("KubeletServiceAccountTokenForCredentialProviders").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateListFromCacheSnapshot = newFeatureGate("ListFromCacheSnapshot").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateLocalStorageCapacityIsolationFSQuotaMonitoring = newFeatureGate("LocalStorageCapacityIsolationFSQuotaMonitoring").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateManifestBasedAdmissionControlConfig = newFeatureGate("ManifestBasedAdmissionControlConfig").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMatchLabelKeysInPodTopologySpread = newFeatureGate("MatchLabelKeysInPodTopologySpread").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMatchLabelKeysInPodTopologySpreadSelectorMerge = newFeatureGate("MatchLabelKeysInPodTopologySpreadSelectorMerge").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMaxUnavailableStatefulSet = newFeatureGate("MaxUnavailableStatefulSet").
						reportProblemsToJiraComponent("apps").
						contactPerson("atiratree").
						productScope(kubernetes).
						enhancementPR("https://github.com/kubernetes/enhancements/issues/961").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMemoryQoS = newFeatureGate("MemoryQoS").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMutablePVNodeAffinity = newFeatureGate("MutablePVNodeAffinity").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMutablePodResourcesForSuspendedJobs = newFeatureGate("MutablePodResourcesForSuspendedJobs").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMutableSchedulingDirectivesForSuspendedJobs = newFeatureGate("MutableSchedulingDirectivesForSuspendedJobs").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateMutatingAdmissionPolicy = newFeatureGate("MutatingAdmissionPolicy").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("benluddy").
						productScope(kubernetes).
						enhancementPR("https://github.com/kubernetes/enhancements/issues/3962").
						enable(inDefault(),inOKD(),inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateNodeDeclaredFeatures = newFeatureGate("NodeDeclaredFeatures").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateNominatedNodeNameForExpectation = newFeatureGate("NominatedNodeNameForExpectation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateOpenAPIEnums = newFeatureGate("OpenAPIEnums").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateOpportunisticBatching = newFeatureGate("OpportunisticBatching").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePLEGOnDemandRelist = newFeatureGate("PLEGOnDemandRelist").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePersistentVolumeClaimUnusedSinceTime = newFeatureGate("PersistentVolumeClaimUnusedSinceTime").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodAndContainerStatsFromCRI = newFeatureGate("PodAndContainerStatsFromCRI").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodCertificateRequest = newFeatureGate("PodCertificateRequest").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade(),withGroupKindResources(groupKindResource{Group: "certificates.k8s.io", Kind: "PodCertificateRequest", Resource: "podcertificaterequests"})).
						mustRegister()


	FeatureGatePodDeletionCost = newFeatureGate("PodDeletionCost").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodLevelResourceManagers = newFeatureGate("PodLevelResourceManagers").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodLevelResources = newFeatureGate("PodLevelResources").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodLogsQuerySplitStreams = newFeatureGate("PodLogsQuerySplitStreams").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodReadyToStartContainersCondition = newFeatureGate("PodReadyToStartContainersCondition").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodTopologyLabelsAdmission = newFeatureGate("PodTopologyLabelsAdmission").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePodsAPI = newFeatureGate("PodsAPI").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePortForwardWebsockets = newFeatureGate("PortForwardWebsockets").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGatePreventStaticPodAPIReferences = newFeatureGate("PreventStaticPodAPIReferences").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateQOSReserved = newFeatureGate("QOSReserved").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateReduceDefaultCrashLoopBackOffDecay = newFeatureGate("ReduceDefaultCrashLoopBackOffDecay").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateRelaxedServiceNameValidation = newFeatureGate("RelaxedServiceNameValidation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateReloadKubeletClientCAFile = newFeatureGate("ReloadKubeletClientCAFile").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateReloadKubeletServerCertificateFile = newFeatureGate("ReloadKubeletServerCertificateFile").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateRemoteRequestHeaderUID = newFeatureGate("RemoteRequestHeaderUID").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateResourceHealthStatus = newFeatureGate("ResourceHealthStatus").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateResourceHealthStatusMessage = newFeatureGate("ResourceHealthStatusMessage").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateRestartAllContainersOnContainerExits = newFeatureGate("RestartAllContainersOnContainerExits").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateRotateKubeletServerCertificate = newFeatureGate("RotateKubeletServerCertificate").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateRuntimeClassInImageCriApi = newFeatureGate("RuntimeClassInImageCriApi").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateSELinuxMount = newFeatureGate("SELinuxMount").
						reportProblemsToJiraComponent("Storage / Kubernetes").
						contactPerson("jsafrane").
						productScope(kubernetes).
						enhancementPR("https://github.com/kubernetes/enhancements/issues/1710").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateSchedulerAsyncAPICalls = newFeatureGate("SchedulerAsyncAPICalls").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateSchedulerAsyncPreemption = newFeatureGate("SchedulerAsyncPreemption").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateSchedulerPopFromBackoffQ = newFeatureGate("SchedulerPopFromBackoffQ").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateServiceAccountNodeAudienceRestriction = newFeatureGate("ServiceAccountNodeAudienceRestriction").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateShardedListAndWatch = newFeatureGate("ShardedListAndWatch").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateSizeBasedListCostEstimate = newFeatureGate("SizeBasedListCostEstimate").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStaleControllerConsistencyDaemonSet = newFeatureGate("StaleControllerConsistencyDaemonSet").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStaleControllerConsistencyJob = newFeatureGate("StaleControllerConsistencyJob").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStaleControllerConsistencyReplicaSet = newFeatureGate("StaleControllerConsistencyReplicaSet").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStaleControllerConsistencyStatefulSet = newFeatureGate("StaleControllerConsistencyStatefulSet").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStatefulSetSemanticRevisionComparison = newFeatureGate("StatefulSetSemanticRevisionComparison").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStorageCapacityScoring = newFeatureGate("StorageCapacityScoring").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStorageVersionAPI = newFeatureGate("StorageVersionAPI").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStorageVersionHash = newFeatureGate("StorageVersionHash").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStorageVersionMigrator = newFeatureGate("StorageVersionMigrator").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStrictIPCIDRValidation = newFeatureGate("StrictIPCIDRValidation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStructuredAuthenticationConfigurationEgressSelector = newFeatureGate("StructuredAuthenticationConfigurationEgressSelector").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateStructuredAuthenticationConfigurationJWKSMetrics = newFeatureGate("StructuredAuthenticationConfigurationJWKSMetrics").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTaintTolerationComparisonOperators = newFeatureGate("TaintTolerationComparisonOperators").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTokenRequestServiceAccountUIDValidation = newFeatureGate("TokenRequestServiceAccountUIDValidation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTopologyAwareWorkloadScheduling = newFeatureGate("TopologyAwareWorkloadScheduling").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTopologyManagerPolicyAlphaOptions = newFeatureGate("TopologyManagerPolicyAlphaOptions").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTopologyManagerPolicyBetaOptions = newFeatureGate("TopologyManagerPolicyBetaOptions").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateTranslateStreamCloseWebsocketRequests = newFeatureGate("TranslateStreamCloseWebsocketRequests").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateUnauthenticatedHTTP2DOSMitigation = newFeatureGate("UnauthenticatedHTTP2DOSMitigation").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateUnknownVersionInteroperabilityProxy = newFeatureGate("UnknownVersionInteroperabilityProxy").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateUnlockWhileProcessingFIFO = newFeatureGate("UnlockWhileProcessingFIFO").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateUserNamespacesHostNetworkSupport = newFeatureGate("UserNamespacesHostNetworkSupport").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateVolumeLimitScaling = newFeatureGate("VolumeLimitScaling").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWatchCacheInitializationPostStartHook = newFeatureGate("WatchCacheInitializationPostStartHook").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWatchList = newFeatureGate("WatchList").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWatchListClient = newFeatureGate("WatchListClient").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWindowsCPUAndMemoryAffinity = newFeatureGate("WindowsCPUAndMemoryAffinity").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWindowsGracefulNodeShutdown = newFeatureGate("WindowsGracefulNodeShutdown").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inTechPreviewNoUpgrade(),inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWorkloadAwarePreemption = newFeatureGate("WorkloadAwarePreemption").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()


	FeatureGateWorkloadWithJob = newFeatureGate("WorkloadWithJob").
						reportProblemsToJiraComponent("kube-apiserver").
						contactPerson("bpalmer").
						productScope(kubernetes).
						enhancementPR("https://github.com/openshift/enhancements/pull/2084").
						enable(inDevPreviewNoUpgrade()).
						mustRegister()

)
