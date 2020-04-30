package clusters

type ClusterInfo struct {
	Cloud         string //GCP, Azure, AWS, or Hub (for a standard neutral format)
	Scope         string //Project in GKE, Subscription in AKS, blank in EKS
	Location      string
	Name          string
	NodeCount     int32
	sourceCluster *ClusterInfo
}

type ClusterAccess interface {
	ListClusters(project, location string) ([]ClusterInfo, error)
	CreateCluster(info ClusterInfo) error
}
