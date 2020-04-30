package cluster_access

import (
	"clusterCloner/clusters/cluster_info"
)

type ClusterAccess interface {
	ListClusters(project, location string) ([]cluster_info.ClusterInfo, error)
	CreateCluster(info cluster_info.ClusterInfo) error
}
