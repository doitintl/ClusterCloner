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

func parseClusterDescription(jsonBytes []byte) ([]EKSCluster, error) {
	eksClusters := make([]EKSCluster, 0)
	err := json.Unmarshal(jsonBytes, &eksClusters)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall "+string(jsonBytes))
	}
	return eksClusters, nil
}
