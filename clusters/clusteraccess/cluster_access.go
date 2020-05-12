package clusteraccess

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clouds/gke/accessgke"
)

// ClusterAccess ...
type ClusterAccess interface {
	//todo allow listing clusters by tag
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]*clusters.ClusterInfo, error)
	//CreateCluster ...
	CreateCluster(info *clusters.ClusterInfo) (*clusters.ClusterInfo, error)
	//GetSupportedK8sVersions
	GetSupportedK8sVersions(scope, location string) []string
}

// GetClusterAccessor ...
func GetClusterAccessor(cloud string) ClusterAccess {
	var clusterAccessor ClusterAccess
	switch cloud {
	case clusters.GCP:
		clusterAccessor = accessgke.GKEClusterAccess{}
	case clusters.AZURE:
		clusterAccessor = accessaks.AKSClusterAccess{}
	default:
		clusterAccessor = nil
	}
	return clusterAccessor
}
