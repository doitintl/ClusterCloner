package cluster_transformation

import (
	"clusterCloner/clusters/clouds/aks/transform"
	transform2 "clusterCloner/clusters/clouds/gke/transform"
	"clusterCloner/clusters/cluster_info"
	"errors"
	"fmt"
	"log"
)

func Transform(clusterInfo cluster_info.ClusterInfo, toCloud string) (c cluster_info.ClusterInfo, err error) {
	hub, err1 := toHubFormat(clusterInfo)
	if err1 != nil {
		return cluster_info.ClusterInfo{}, err1
	}
	target, err2 := fromHubFormat(hub, toCloud)
	if err2 != nil {
		return cluster_info.ClusterInfo{}, err2
	}
	return target, nil
}
func toHubFormat(clusterInfo cluster_info.ClusterInfo) (c cluster_info.ClusterInfo, err error) {
	err = nil
	var ret cluster_info.ClusterInfo
	switch cloud := clusterInfo.Cloud; cloud { //todo Split out into "adapters" to avoid switch statement. Putting it here now for reference.
	case cluster_info.GCP:
		log.Print("From GCP ")
		ret, err = transform2.TranformGCPToHub(clusterInfo)
	case cluster_info.AZURE:
		log.Print("From Azure ")
		ret, err = transform.TransformAzureToHub(clusterInfo)
	case cluster_info.AWS:
		log.Print("From AWS ")
		return c, errors.New(fmt.Sprintf("Unsupported %s", cloud))
	case cluster_info.HUB:
		log.Print("From Hub , no changes")
		ret = c
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", cloud))
	}
	return ret, err
}
func fromHubFormat(clusterInfo cluster_info.ClusterInfo, toCloud string) (c cluster_info.ClusterInfo, err error) {
	if clusterInfo.Cloud != cluster_info.HUB {
		return cluster_info.ClusterInfo{}, errors.New(fmt.Sprintf("Wrong Cloud %s", clusterInfo.Cloud))
	}
	err = nil
	var ret cluster_info.ClusterInfo
	switch toCloud { //todo Split out into "adapters" to avoid switch statement. Putting it here now for reference.
	case cluster_info.GCP:
		log.Print("to GCP ")
		ret, err = transform2.TransformHubToGCP(clusterInfo)
	case cluster_info.AZURE:
		log.Print("to Azure ")
		ret, err = transform.TransformHubToAzure(clusterInfo)
	case cluster_info.AWS:
		log.Print("to AWS ")
		return c, errors.New(fmt.Sprintf("Unsupported %s", toCloud))
	case cluster_info.HUB:
		log.Print("to Hub , no changes")
		ret = c
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", toCloud))
	}
	return ret, err
}
