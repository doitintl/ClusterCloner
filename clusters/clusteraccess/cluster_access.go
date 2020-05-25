package clusteraccess

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	accesseks "clustercloner/clusters/clouds/eks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
)

// ClusterAccess ...
type ClusterAccess interface {
	//List list all clusters at this location
	List(project, location string, labels map[string]string) ([]*clusters.ClusterInfo, error)
	//Create ...
	Create(info *clusters.ClusterInfo) (*clusters.ClusterInfo, error)
	//Describe...
	Describe(readThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error)
	//Delete
	Delete(ci *clusters.ClusterInfo) error
	//GetSupportedK8sVersions
	GetSupportedK8sVersions(scope, location string) []string
}

// GetClusterAccess ...
func GetClusterAccess(cloud string) ClusterAccess {
	var clusterAccessor ClusterAccess
	switch cloud {
	case clusters.GCP:
		clusterAccessor = accessgke.GKEClusterAccess{}
	case clusters.Azure:
		clusterAccessor = accessaks.AKSClusterAccess{}
	case clusters.AWS:
		clusterAccessor = accesseks.EKSClusterAccess{}
	default:
		panic("unsupported cloud " + cloud)
	}
	return clusterAccessor
}
