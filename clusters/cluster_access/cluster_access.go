package cluster_access

import (
	"clusterCloner/clusters/cluster_info"
)

type ClusterAccess interface {
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]cluster_info.ClusterInfo, error)
	CreateCluster(info cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error)
}
type Transformer interface {
	// CloudToHub todo: Extract CloudToHub and HubToCloud as 'embedded' functions to be shared by implementors
	CloudToHub(in cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error)
	//	HubToCloud///
	HubToCloud(in cluster_info.ClusterInfo, outputScope string) (cluster_info.ClusterInfo, error)
	// LocationHubToCloud ...
	LocationHubToCloud(loc string) (string, error)
	// LocationCloudToHub ...
	LocationCloudToHub(loc string) (string, error)
}
