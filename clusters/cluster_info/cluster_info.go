package cluster_info

type ClusterInfo struct {
	Cloud         string //GCP, Azure, AWS, or Hub (for a standard neutral format)
	Scope         string //Project in GKE, Subscription in AKS, blank in EKS
	Location      string
	Name          string
	K8sVersion    string
	NodeCount     int32
	GeneratedBy   string
	SourceCluster *ClusterInfo
}

var (
	MOCK           = "Mock"
	READ           = "Read"
	CREATED        = "Created"
	TRANSFORMATION = "Transformation"
)

var (
	HUB   = "Hub"
	GCP   = "GCP"
	AZURE = "Azure"
	AWS   = "AWS"
)
