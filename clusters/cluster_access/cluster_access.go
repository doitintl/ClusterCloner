package cluster_access

import (
	"clusterCloner/clusters/cluster_info"
)

type ClusterAccess interface {
	ListClusters(project, location string) ([]cluster_info.ClusterInfo, error)
	CreateCluster(info cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error)
}
type Transformer interface {
	CloudToHub(clusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error)
	HubToCloud(clusterInfo cluster_info.ClusterInfo, outputScope string) (cluster_info.ClusterInfo, error)
	LocationHubToCloud(loc string) (string, error)
	LocationCloudToHub(loc string) (string, error)
}
