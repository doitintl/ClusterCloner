package clusters

import (
	"errors"
	"fmt"
	"log"
)

func ToHubFormat(clusterInfo ClusterInfo) (c ClusterInfo, err error) {
	err = nil
	var ret ClusterInfo
	switch cloud := clusterInfo.Cloud; cloud { //todo Split out into "adapters" to avoid switch statement. Putting it here now for reference.
	case GCP:
		log.Print("From GCP to Hub")
		ret, err = ConvertGCPToStandard(clusterInfo)
	case AZURE:
		log.Print("From Azure to Hub")
		ret, err = ConvertAzureToStandard(clusterInfo)
	case AWS:
		log.Print("From AWS to Hub")
		return c, errors.New(fmt.Sprintf("Unsupported %s", cloud))
	case HUB:
		log.Print("From Hub to Hub, no changes")
		ret = c
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", cloud))
	}
	return ret, err
}

func ConvertAzureToStandard(clusterInfo ClusterInfo) (ClusterInfo, error) {
	var ret = clusterInfo
	ret.Cloud = HUB
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in conversion cross-cloud
	loc, err := ConvertLocationAzureToHub(ret.Location)
	ret.Location = loc
	return ret, err
}

func ConvertGCPToStandard(clusterInfo ClusterInfo) (ClusterInfo, error) {
	var ret = clusterInfo
	ret.Cloud = HUB
	//	ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in converstion cross-cloud
	loc, err := ConvertLocationGcpToHub(ret.Location)
	ret.Location = loc
	return ret, err
}
