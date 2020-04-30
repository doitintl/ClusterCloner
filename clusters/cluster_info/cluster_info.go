package cluster_info

type ClusterInfo struct {
	Cloud         string //GCP, Azure, AWS, or Hub (for a standard neutral format)
	Scope         string //Project in GKE, Subscription in AKS, blank in EKS
	Location      string
	Name          string
	NodeCount     int32
	SourceCluster *ClusterInfo
}

var (
	HUB   = "Hub"
	GCP   = "GCP"
	AZURE = "Azure"
	AWS   = "AWS"
)
