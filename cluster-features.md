# Cluster Features

## Key
- "Yes": features that are copied
- "No": features that will not be copied.
  - If not marked, it is because the feature is cloud-specific.
  - Other features are not copied becuase they are ephemeral, like cluster status, or secondary -- calculated from other values.
- "TBD": features that should be copied

## Fields of Cluster (Mostly based on GKE)
- Name: Yes
- Description: No
- InitialNodeCount: No (Use value in NodePool)
- NodeConfig: No (Use value in NodePool)
- MasterAuth: No
- LoggingService: No
- MonitoringService: No
- Network: TBD (Convert this in AWS to VPC; need to create system for re-use of VPCs or creation of VPCs)
- ClusterIpv4Cidr: TBD
- AddonsConfig: No
- Subnetwork TBD (See above re AWS VPC.)
- NodePools: Yes
- Locations: Yes
- EnableKubernetesAlpha: No
- ResourceLabels: Yes
- LabelFingerprint: No
- LegacyAbac: No
- NetworkPolicy: TBD
- IpAllocationPolicy: No
- MasterAuthorizedNetworksConfig: No
- MaintenancePolicy: No (Maintenance window)
- NetworkConfig: No
- PrivateClusterConfig: TBD
- SelfLink: No (secondary value)
- Zone: TBD in full. Allow reading and specification of EKS and
AKS availability zones. As-is, GKE zonal regional clusters can be read,
but the user has to specify which zone.)
- Endpoint: No
- InitialClusterVersion: Yes (from CurrentMasterVersion)
- CurrentMasterVersion: Yes  (into Initial Cluster Version)
- CurrentNodeVersion: No
- CreateTime: No (ephemeral value)
- Status: No (ephemeral value)
- StatusMessage: No (ephemeral value)
- NodeIpv4CidrSize: No
- ServicesIpv4Cidr: No
- InstanceGroupUrls: No (secondary value)
- CurrentNodeCount: No (ephemeral)
- ExpireTime: No (ephemeral)
- Location: Yes
- Unmanaged clusters:  No
-----------------------------
## Fields of Node (Based on GKE)
- Name: Yes
- Config: See below
- InitialNodeCount: Yes
- SelfLink: No (secondary value)
- Version: TBD (now, we just re-use cluster K8s version here)
- InstanceGroupUrl: No (secondary value)
- Status: No (ephemeral)
- StatusMessage: No (ephemeral)
- Autoscaling: TBD
- Management: TBD (Auto-repair and Auto-upgrade)
- InstanceType: Support multi-instance type
-------------------------------
## Fields of NodeConfig (Based on GKE)
-  MachineType: Yes
-  DiskSizeGb: Yes
-  OauthScopes: No
-  ServiceAccount: No
-  Metadata: No
-  ImageType: No
-  Labels: TBD
-  LocalSsdCount: No (Maybe should)
-  Tags: TBD (as-is, only cluster tags are copied)
-  Preemptible: Partly (GKE and ASK only)
-  Accelerators: TBD (In EKS, this is defined differently: use a different machine type to get a GPU.)
-  DiskType: TBD

### Clusters and NodeGroups
- Use Paging in EKS and AKS

### Converting reference data
#### Machine Types
Now converted with a simple algorithm: We choose the smallest target machine type bigger than the input machine type in CPU and RAM.

The results of the algorithm may be low-quality. Improve this algorithm, or created a manual conversion table.

### Kubernetes versions
Now converted with a simple algorithm:

- For clouds that support patch versions, namely AKS and GKE, we choose the least patch version that is
- no less than the supplied version, but has the same minor version.
- If that is not possible, we get the largest patch version that has the same minor version.
- For EKS, which does not support patch versions, we just choose the same minor version.

The results of the algorithm may be low-quality. Improve it.