package clusteraccess

import (
	"clustercloner/clusters/clusterinfo"
)

// ClusterAccess ...
type ClusterAccess interface {
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]clusterinfo.ClusterInfo, error)
	CreateCluster(info clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error)
}

// Transformer ...
type Transformer interface {
	// CloudToHub todo: Extract CloudToHub and HubToCloud as 'embedded' functions to be shared by implementors
	CloudToHub(in clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error)
	//	HubToCloud///
	HubToCloud(in clusterinfo.ClusterInfo, outputScope string) (clusterinfo.ClusterInfo, error)
	// LocationHubToCloud ...
	LocationHubToCloud(loc string) (string, error)
	// LocationCloudToHub ...
	LocationCloudToHub(loc string) (string, error)
}
