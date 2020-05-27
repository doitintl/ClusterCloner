package eksctl

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// EKSCluster ...
type EKSCluster struct {
	Name    string
	Tags    map[string]string
	Version string
}

// EKSNodeGroup ...
type EKSNodeGroup struct {
	Name string
	//TODO, need to support autoscaling in NodeGroupsin all Clouds. Until then, we read this from EKS to stand-in for size of static NodeGroup
	DesiredCapacity int
	InstanceType    string
}

func parseClusterDescription(jsonBytes []byte) ([]EKSCluster, error) {
	eksClusters := make([]EKSCluster, 0)
	err := json.Unmarshal(jsonBytes, &eksClusters)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall "+string(jsonBytes))
	}
	return eksClusters, nil
}

func parseNodeGroupsDescription(jsonBytes []byte) ([]EKSNodeGroup, error) {
	eksNodeGroups := make([]EKSNodeGroup, 0)
	err := json.Unmarshal(jsonBytes, &eksNodeGroups)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall "+string(jsonBytes))
	}
	return eksNodeGroups, nil
}
