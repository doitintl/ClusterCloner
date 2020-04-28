package clusters

type ClusterInfo struct {
	Name      string
	NodeCount int32
}

type ClusterAccess interface {
	ListClusters(project, location string) ([]ClusterInfo, error)
}
