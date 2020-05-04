package cluster_transformation

import (
	accessaks "clusterCloner/clusters/clouds/aks/access"
	transformaks "clusterCloner/clusters/clouds/aks/transform"
	accesseks "clusterCloner/clusters/clouds/eks/access"
	accessgke "clusterCloner/clusters/clouds/gke/access"
	transformgke "clusterCloner/clusters/clouds/gke/transform"
	"clusterCloner/clusters/cluster_access"
	"clusterCloner/clusters/cluster_info"
	"fmt"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	"log"
)

func Clone(cliCtx *cli.Context) ([]cluster_info.ClusterInfo, error) {
	inputCloud := cliCtx.String("inputcloud")
	outputCloud := cliCtx.String("outputcloud")
	inputLocation := cliCtx.String("inputlocation")
	inputScope := cliCtx.String("inputscope")
	outputScope := cliCtx.String("outputscope")
	create := cliCtx.Bool("create")
	sCreate := ""
	if !create {
		sCreate = "not "
	}
	log.Printf("Will %screate target clusters", sCreate)
	if inputCloud == "" || outputCloud == "" || inputLocation == "" || inputScope == "" || outputScope == "" {
		log.Fatal("Missing flags")
	}
	return clone(inputCloud, outputCloud, inputLocation, inputScope, outputScope, create)

}

func clone(inputCloud string, outputCloud string, inputLocation string, inputScope string, outputScope string, create bool) ([]cluster_info.ClusterInfo, error) {
	var clusterAccessor cluster_access.ClusterAccess
	switch inputCloud {
	case cluster_info.GCP:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (project) for GCP")
		}
		clusterAccessor = accessgke.GkeClusterAccess{}
	case cluster_info.AZURE:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (Resource Group) for Azure")
		}
		clusterAccessor = accessaks.AksClusterAccess{}
	case cluster_info.AWS:
		clusterAccessor = accesseks.EksClusterAccess{}
		//todo support CloudToHub as an output cloud
	default:
		return nil, errors.New(fmt.Sprintf("Invalid inputcloud \"%s\"", inputCloud))
	}

	inputClusterInfos, err := clusterAccessor.ListClusters(inputScope, inputLocation)
	if err != nil {
		return nil, err
	}
	outputClusterInfos := make([]cluster_info.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := transform(inputClusterInfo, outputCloud, outputScope)
		outputClusterInfos = append(outputClusterInfos, outputClusterInfo)
		if err != nil {
			log.Printf("Error processing %v: %v", inputClusterInfo, err)
		}
	}
	createdClusterInfos := make([]cluster_info.ClusterInfo, len(outputClusterInfos))

	if create {
		log.Println("Creating", len(outputClusterInfos), "target clusters")
		for idx, createThis := range outputClusterInfos {
			createCi, err_ := CreateCluster(createThis)
			if err_ != nil {
				log.Printf("Error creating %v: %v", createThis, err)
			}
			createdClusterInfos[idx] = createCi
		}
		blank := cluster_info.ClusterInfo{}
		for idx, createdClusterInfo := range createdClusterInfos {
			if createdClusterInfo != blank {
				outputClusterInfos[idx] = createdClusterInfo
			}
		}
		if len(outputClusterInfos) != len(inputClusterInfos) {
			log.Fatalf("%d != %d", len(outputClusterInfos), len(inputClusterInfos))
		}
		return outputClusterInfos, nil //replaced each ClusterInfo that was created; left the ones that were not
	} else {
		log.Println("Not creating", len(outputClusterInfos), "target clusters")
		return outputClusterInfos, nil
	}
}

func CreateCluster(createThis cluster_info.ClusterInfo) (createdClusterInfo cluster_info.ClusterInfo, err error) {
	var ca cluster_access.ClusterAccess
	switch createThis.Cloud {
	case cluster_info.AZURE:
		ca = accessaks.AksClusterAccess{}
	case cluster_info.AWS:
		panic("AWS not implemented")
	case cluster_info.GCP:
		ca = accessgke.GkeClusterAccess{}
	default:
		return cluster_info.ClusterInfo{}, errors.New("Unsupported Cloud for creating a cluster: " + createThis.Cloud)

	}
	created, err := ca.CreateCluster(createThis)
	if err != nil {
		log.Println("Error creating cluster", err)
	}
	return created, err
}

func transform(clusterInfo cluster_info.ClusterInfo, toCloud string, outputScope string) (c cluster_info.ClusterInfo, err error) {
	if clusterInfo.Cloud == "" {
		return c, errors.New("No cloud name found in input cluster")
	}
	hub, err1 := toHubFormat(clusterInfo)
	if err1 != nil {
		return cluster_info.ClusterInfo{}, errors.Wrap(err1, "Error in transforming to CloudToHub Format")
	}
	target, err2 := fromHubFormat(hub, toCloud, outputScope)
	if err2 != nil {
		return cluster_info.ClusterInfo{}, err2
	}
	if clusterInfo.Cloud == toCloud {
		// todo maybe self-to-self transformation shoud not go thru hub format.
		target.Name = target.Name + "-copy"
	}
	target.GeneratedBy = cluster_info.TRANSFORMATION
	return target, nil
}
func toHubFormat(input cluster_info.ClusterInfo) (c cluster_info.ClusterInfo, err error) {
	err = nil
	var ret cluster_info.ClusterInfo
	var transformer cluster_access.Transformer
	switch cloud := input.Cloud; cloud {
	case cluster_info.GCP:
		transformer = transformgke.GkeTransformer{}
	case cluster_info.AZURE:
		transformer = transformaks.AksTransformer{}
	case cluster_info.AWS:
		return c, errors.New(fmt.Sprintf("Unsupported %s", cloud))
	case cluster_info.HUB:
		log.Print("From CloudToHub , no changes")
		ret = input
		return ret, nil
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", cloud))
	}
	ret, err = transformer.CloudToHub(input)
	return ret, err
}
func fromHubFormat(hub cluster_info.ClusterInfo, toCloud string, outputScope string) (c cluster_info.ClusterInfo, err error) {
	if hub.Cloud != cluster_info.HUB {
		return cluster_info.ClusterInfo{}, errors.New(fmt.Sprintf("Wrong Cloud %s", hub.Cloud))
	}
	var transformer cluster_access.Transformer
	err = nil
	var ret cluster_info.ClusterInfo
	switch toCloud { //  We do not expect more than  these clouds so not splitting out dynamically loaded adapters
	case cluster_info.GCP:
		transformer = transformgke.GkeTransformer{}
	case cluster_info.AZURE:
		transformer = transformaks.AksTransformer{}
	case cluster_info.AWS:
		return c, errors.New(fmt.Sprintf("Unsupported %s", toCloud))
	case cluster_info.HUB:
		log.Print("to Hub , no changes")
		ret = hub
		return ret, nil
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", toCloud))
	}
	ret, err = transformer.HubToCloud(hub, outputScope)
	return ret, err
}
