package clusters

// ClusterInfo ...
type ClusterInfo struct {
	Cloud         string //GCP, Azure, AWS, or Hub (for a standard neutral format)
	Scope         string //Project in GKE, Subscription in AKS, blank in EKS
	Location      string //Zone or Region in GKE, Region in others
	Name          string
	K8sVersion    string
	GeneratedBy   string
	NodePools     []NodePoolInfo
	SourceCluster *ClusterInfo
}

// AddNodePool ...
func (ci *ClusterInfo) AddNodePool(pool NodePoolInfo) {
	if ci.NodePools == nil {
		ci.NodePools = make([]NodePoolInfo, 0)
	}
	ci.NodePools = append(ci.NodePools, pool)
}

// MachineType ...
type MachineType struct {
	Name  string
	CPU   int32
	RAMMB int32
}

// NodePoolInfo ...
type NodePoolInfo struct {
	Name        string
	NodeCount   int32
	K8sVersion  string
	MachineType MachineType
	DiskSizeGB  int32
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
