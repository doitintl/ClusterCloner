package clusters

// ClusterAccess ...
type ClusterAccess interface {
	//todo allow listing clusters by tag
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]*ClusterInfo, error)
	//CreateCluster ...
	CreateCluster(info *ClusterInfo) (*ClusterInfo, error)
}
