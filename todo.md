
# Plans to copy more fields

## Fields of Cluster (Based on GKE)
- Name: Copied
- Description: Will not copy (Cloud-specific)
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
## Fields of Node (Based on GKE)
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
- InstanceType: Support multi-instance type
-------------------------------
## Feelds of NodeConfig (Based on GKE)
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


### Clusters and NodeGroups
- Use Paging in EKS and AKS

### Converting reference data
#### Machine Types
Now converted with a simple algorithm. (We choose the smallest target machine type bigger than the input machine type in CPU and RAM.)
The results of the algorithm may be low-quality. Improve this algorithm, or created a manual conversion table.
### Kubernetes versions
Now converted with a simple algorithm. (For clouds that support patch versions, namely AKS and GKE,
we choose the least patch version that is
greater or equal to  the supplied version, but has the same major-minor version.
If that is not possible, we get the largest patch version that has the same major-minor version.
For EKS, that does not support patch versions, we just choose the same major-minor version.
The results of the algorithm may be low-quality. Improve this algorithm.