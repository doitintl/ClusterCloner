package cluster_transformation

import (
	access_aks "clusterCloner/clusters/clouds/aks/access"
	transform_aks "clusterCloner/clusters/clouds/aks/transform"
	access_eks "clusterCloner/clusters/clouds/eks/access"
	access_gke "clusterCloner/clusters/clouds/gke/access"
	transform_gke "clusterCloner/clusters/clouds/gke/transform"
	"clusterCloner/clusters/cluster_access"
	"clusterCloner/clusters/cluster_info"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
)

func Clone(ctx *cli.Context) ([]cluster_info.ClusterInfo, error) {
	var cluster_accessor cluster_access.ClusterAccess
	switch ctx.String("inputcloud") {
	case cluster_info.GCP:
		cluster_accessor = access_gke.GkeClusterAccess{}
	case cluster_info.AZURE:
		cluster_accessor = access_aks.AksClusterAccess{}
	case cluster_info.AWS:
		cluster_accessor = access_eks.EksClusterAccess{}

	}

	inputClusterInfos, err := cluster_accessor.ListClusters("", "")
	if err != nil {
		return nil, err
	}
	outputClusterInfos := make([]cluster_info.ClusterInfo, 10)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := Transform(inputClusterInfo, ctx.String("outputcloud"))
		outputClusterInfos = append(outputClusterInfos, outputClusterInfo)
		if err != nil {
			log.Printf("Error processing %v: %v", inputClusterInfo, err)
		}
	}
	createdClusterInfos := make([]cluster_info.ClusterInfo, len(outputClusterInfos))

	if !ctx.Bool("create") {
		for idx, createThis := range outputClusterInfos {
			createCi, err_ := CreateCluster(createThis)
			if err_ != nil {
				log.Printf("Error creating %v: %v", createThis, err)
			}
			createdClusterInfos[idx] = createCi
		}
	}
	blank := cluster_info.ClusterInfo{}
	for idx, createdClusterInfo := range createdClusterInfos {
		if createdClusterInfo != blank {
			outputClusterInfos[idx] = createdClusterInfo
		}
	}
	return outputClusterInfos, nil

}

func CreateCluster(createThis cluster_info.ClusterInfo) (createdClusterInfo cluster_info.ClusterInfo, err error) {
	panic("implement")

	createdClusterInfo.GeneratedBy = cluster_info.CREATED
	return cluster_info.ClusterInfo{}, nil
}

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
		ret, err = transform_gke.TranformGCPToHub(clusterInfo)
	case cluster_info.AZURE:
		log.Print("From Azure ")
		ret, err = transform_aks.TransformAzureToHub(clusterInfo)
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
		ret, err = transform_gke.TransformHubToGCP(clusterInfo)
	case cluster_info.AZURE:
		log.Print("to Azure ")
		ret, err = transform_aks.TransformHubToAzure(clusterInfo)
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
