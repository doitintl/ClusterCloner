package clusteraccess

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clouds/gke/access"
	"log"
)

// ClusterAccess ...
type ClusterAccess interface {
	//TODO allow listing clusters by tag
	//ListClusters list all clusters at this location
	ListClusters(project, location string) ([]*clusters.ClusterInfo, error)
	//CreateCluster ...
	CreateCluster(info *clusters.ClusterInfo) (*clusters.ClusterInfo, error)
	//GetSupportedK8sVersions
	GetSupportedK8sVersions(scope, location string) []string
}

// GetClusterAccess ...
func GetClusterAccess(cloud string) ClusterAccess {
	var clusterAccessor ClusterAccess
	switch cloud {
	case clusters.GCP:
		clusterAccessor = access.GKEClusterAccess{}
	case clusters.AZURE:
		clusterAccessor = accessaks.AKSClusterAccess{}
	default:
		log.Println("unsupported cloud ", cloud)
		clusterAccessor = nil
	}
	return clusterAccessor
}
