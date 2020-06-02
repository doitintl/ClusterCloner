package awssdk

/*
// EKSNodeGroup ...
type EKSNodeGroup struct {
	ClusterName   string
	DiskSize      int
	InstanceTypes []string
	ScalingConfig map[string]int
	Labels        map[string]string
	Tags          map[string]string
	Version       string
}

// NGHolder ...
type NGHolder struct {
	NodeGroup EKSNodeGroup
}

func parseNodeGroupDescription(jsonBytes []byte) (*EKSNodeGroup, error) {
	eksNodeGroupHolder := &NGHolder{}
	err := json.Unmarshal(jsonBytes, &eksNodeGroupHolder)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall "+string(jsonBytes))
	}
	return &eksNodeGroupHolder.NodeGroup, nil
}
*/
