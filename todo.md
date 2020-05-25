
# Plans to copy more fields

## Fields of GKE Cluster
- Name: Copied
- Description: Should Copy
- InitialNodeCount: Will not copy (in NodePool)
- NodeConfig: Will not copy (in NodePool)
- MasterAuth: Will not copy (Cloud-specific)
- LoggingService: Will not copy (Cloud-specific)
- MonitoringService: Will not copy (Cloud-specific)
- Network: Should copy (to AWS VPC; need to create system for re-use of VPCs or creation of VPCs)
- ClusterIpv4Cidr: Should copy
- AddonsConfig: Will not copy (Cloud-specific)
- Subnetwork Copy (to AWS VPC; need to create system for re-use of subnetworks or creation of subnetworks)
- NodePools: Copied
- Locations: Copied
- EnableKubernetesAlpha: Need to copy
- ResourceLabels: Copied
- LabelFingerprint: Will not copy (Cloud-specific)
- LegacyAbac: Will not copy (Cloud-specific)
- NetworkPolicy: Should copy
- IpAllocationPolicy: Will not copy (Cloud-specific)
- MasterAuthorizedNetworksConfig: Will not copy (Cloud-specific)
- MaintenancePolicy: Will not copy (Maintenance window; Cloud-specific)
- NetworkConfig: Will not copy
- PrivateClusterConfig: Should copy
- SelfLink: Will not copy (secondary value)
- Zone: Copied (in Location)
- Endpoint: Will not copy
- InitialClusterVersion: Copied (from CurrentMasterVersion)
- CurrentMasterVersion: Copied  (into Initial Cluster Version)
- CurrentNodeVersion: Will not copy
- CreateTime: Will not copy (ephemeral value)
- Status: Will not copy (ephemeral value)
- StatusMessage: Will not copy (Cloud-specific)
- NodeIpv4CidrSize: Will not copy (Cloud-specific)
- ServicesIpv4Cidr: Will not copy (Cloud-specific)
- InstanceGroupUrls: Will not copy (secondary value)
- CurrentNodeCount: Will not copy (ephemeral)
- ExpireTime: Will not copy (ephemeral)
- Location: Copied

-----------------------------
## GKE  Node fields
- Name: Copied
- Config: See below
- InitialNodeCount: Copied
- SelfLink: Will not copy (secondary value)
- Version: Should copy (now, we just re-use cluster K8s version here)
- InstanceGroupUrl: Will not copy (secondary value)
- Status: Will not copy (ephemeral)
- StatusMessage: Will not copy (ephemeral)
- Autoscaling: Should copy
- Management: Should copy (Auto-repair and Auto-upgrade)

## GKE  NodeConfig fields
-  MachineType: Copied
-  DiskSizeGb: Copied
-  OauthScopes: Will not copy (Cloud-specific)
-  ServiceAccount: Will not copy (Cloud-specific)
-  Metadata: Will not copy (Cloud-specific)
-  ImageType: Will not copy (Cloud-specific)
-  Labels: Should copy
-  LocalSsdCount: Will not copy (Maybe should)
-  Tags: Should copy (as-is, only cluster tags are copied)
-  Preemptible: Should copy
-  Accelerators: Should copy (in EKS, to get a GPU, use a different machine type)
-  DiskType: Should copy
-  MinCpuPlatform: Will not copy (Cloud-specific)
