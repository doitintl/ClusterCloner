package clusterinfo

// ClusterInfo ...
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
	// MOCK ...
	MOCK = "Mock"
	// READ ...
	READ = "Read"
	// CREATED ...
	CREATED = "Created"
	// TRANSFORMATION ...
	TRANSFORMATION = "Transformation"
)

var (
	// HUB ...
	HUB = "Hub"
	// GCP ...
	GCP = "GCP"
	// AZURE ...
	AZURE = "Azure"
	// AWS ...
	AWS = "AWS"
)
