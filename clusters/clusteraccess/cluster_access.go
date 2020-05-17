package clusteraccess

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clouds/gke/access"
	"log"
)

// ClusterAccess ...
type ClusterAccess interface {
	//ListClusters list all clusters at this location
	ListClusters(project, location string, labels map[string]string) ([]*clusters.ClusterInfo, error)
	//CreateCluster ...
	CreateCluster(info *clusters.ClusterInfo) (*clusters.ClusterInfo, error)
	//DescribeCluster...
	DescribeCluster(readThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error)
	//GetSupportedK8sVersions
	GetSupportedK8sVersions(scope, location string) []string
}

// GetClusterAccess ...
func GetClusterAccess(cloud string) ClusterAccess {
	var clusterAccessor ClusterAccess
	switch cloud {
	case clusters.GCP:
		clusterAccessor = access.GKEClusterAccess{}
	case clusters.Azure:
		clusterAccessor = accessaks.AKSClusterAccess{}
	default: //TODO support Amazon
		log.Println("unsupported cloud ", cloud)
		clusterAccessor = nil
	}
	return clusterAccessor
}
