package eksctl

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"strings"
)

// EKSCluster ...
type EKSCluster struct {
	Name    string
	Tags    map[string]string
	Version string
}

// EKSNodeGroup ...
type EKSNodeGroup struct { //TODO delete
	Name string
	//TODO, need to support autoscaling in NodeGroups (in all Clouds). This info is not available in eksctl output.
	// Until then, we read this from EKS to stand-in for size of static NodeGroup.
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

//example	"NAME\t\tREGION\nclus-sudic\tus-east-2\n"
func parseClusterList(s string, expectRegion string) ([]string, error) {
	ret := make([]string, 0)
	if strings.Contains(s, "No clusters found") {
		log.Println("Listing clusters: " + s)
		return ret, nil
	}

	sNormalized := strings.ReplaceAll(s, "\t\t", "\t")
	lines := strings.Split(sNormalized, "\n")
	for i, line := range lines {
		parts := strings.Split(line, "\t")
		if line == "" {
			continue
		}
		if len(parts) != 2 {
			return nil, errors.New("wrong number of fields  " + line)
		}
		if i == 0 {
			if line != "NAME\tREGION" {
				return nil, errors.New("bad header line " + line)
			}
			continue
		}

		region := parts[1]
		if region != expectRegion {
			return nil, errors.New("unexpected region " + region + " instead of " + expectRegion)
		}
		clusterName := parts[0]
		ret = append(ret, clusterName)
	}

	return ret, nil
}
