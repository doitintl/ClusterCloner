package clusteraccess

import (
	"clustercloner/clusters/clusterinfo"
)

// ClusterAccess ...
type ClusterAccess interface {
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]*clusterinfo.ClusterInfo, error)
	//CreateCluster ...
	CreateCluster(info *clusterinfo.ClusterInfo) (*clusterinfo.ClusterInfo, error)
}
