package transformation

import (
	accessaks "clustercloner/clusters/clouds/aks/access"
	transformaks "clustercloner/clusters/clouds/aks/transform"
	accesseks "clustercloner/clusters/clouds/eks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
	transformgke "clustercloner/clusters/clouds/gke/transform"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/clusterinfo"
	"fmt"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	"log"
)

// Clone ...
func Clone(cliCtx *cli.Context) ([]clusterinfo.ClusterInfo, error) {
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

func clone(inputCloud string, outputCloud string, inputLocation string, inputScope string, outputScope string, create bool) ([]clusterinfo.ClusterInfo, error) {
	var clusterAccessor clusteraccess.ClusterAccess
	switch inputCloud {
	case clusterinfo.GCP:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (project) for GCP")
		}
		clusterAccessor = accessgke.GkeClusterAccess{}
	case clusterinfo.AZURE:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (Resource Group) for Azure")
		}
		clusterAccessor = accessaks.AksClusterAccess{}
	case clusterinfo.AWS:
		clusterAccessor = accesseks.EksClusterAccess{}
		//todo support CloudToHub as an output cloud
	default:
		return nil, errors.New(fmt.Sprintf("Invalid inputcloud \"%s\"", inputCloud))
	}

	inputClusterInfos, err := clusterAccessor.ListClusters(inputScope, inputLocation)
	if err != nil {
		return nil, err
	}
	outputClusterInfos := make([]clusterinfo.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := transformCloudToCloud(inputClusterInfo, outputCloud, outputScope)
		outputClusterInfos = append(outputClusterInfos, outputClusterInfo)
		if err != nil {
			log.Printf("Error processing %v: %v", inputClusterInfo, err)
		}
	}
	createdClusterInfos := make([]clusterinfo.ClusterInfo, len(outputClusterInfos))

	if create {
		log.Println("Creating", len(outputClusterInfos), "target clusters")
		for idx, createThis := range outputClusterInfos {
			createCi, err := CreateCluster(createThis)
			if err != nil {
				log.Printf("Error creating %v: %v", createThis, err)
			}
			createdClusterInfos[idx] = createCi
		}
		blank := clusterinfo.ClusterInfo{}
		for idx, createdClusterInfo := range createdClusterInfos {
			if createdClusterInfo != blank {
				outputClusterInfos[idx] = createdClusterInfo
			}
		}
		if len(outputClusterInfos) != len(inputClusterInfos) {
			log.Fatalf("%d != %d", len(outputClusterInfos), len(inputClusterInfos))
		}
		return outputClusterInfos, nil //replaced each ClusterInfo that was created; left the ones that were not
	}
	log.Println("Not creating", len(outputClusterInfos), "target clusters")
	return outputClusterInfos, nil

}

// CreateCluster ...
func CreateCluster(createThis clusterinfo.ClusterInfo) (createdClusterInfo clusterinfo.ClusterInfo, err error) {
	var ca clusteraccess.ClusterAccess
	switch createThis.Cloud {
	case clusterinfo.AZURE:
		ca = accessaks.AksClusterAccess{}
	case clusterinfo.AWS:
		panic("AWS not implemented")
	case clusterinfo.GCP:
		ca = accessgke.GkeClusterAccess{}
	default:
		return clusterinfo.ClusterInfo{}, errors.New("Unsupported Cloud for creating a cluster: " + createThis.Cloud)

	}
	created, err := ca.CreateCluster(createThis)
	if err != nil {
		log.Println("Error creating cluster", err)
	}
	return created, err
}

func transformCloudToCloud(clusterInfo clusterinfo.ClusterInfo, toCloud string, outputScope string) (c clusterinfo.ClusterInfo, err error) {
	if clusterInfo.Cloud == "" {
		return c, errors.New("No cloud name found in input cluster")
	}
	hub, err1 := toHubFormat(clusterInfo)
	if err1 != nil {
		return clusterinfo.ClusterInfo{}, errors.Wrap(err1, "Error in transforming to CloudToHub Format")
	}
	target, err2 := fromHubFormat(hub, toCloud, outputScope)
	if err2 != nil {
		return clusterinfo.ClusterInfo{}, err2
	}
	if clusterInfo.Cloud == toCloud {
		// Maybe self-to-self transformation shoud not go thru hub format.
		target.Name = target.Name + "-copy"
	}
	target.GeneratedBy = clusterinfo.TRANSFORMATION
	return target, nil
}
func toHubFormat(input clusterinfo.ClusterInfo) (c clusterinfo.ClusterInfo, err error) {
	err = nil
	var ret clusterinfo.ClusterInfo
	var transformer clusteraccess.Transformer
	switch cloud := input.Cloud; cloud {
	case clusterinfo.GCP:
		transformer = &transformgke.GkeTransformer{}
	case clusterinfo.AZURE:
		transformer = &transformaks.AksTransformer{}
	case clusterinfo.AWS:
		return c, errors.New(fmt.Sprintf("Unsupported %s", cloud))
	case clusterinfo.HUB:
		log.Print("From CloudToHub , no changes")
		ret = input
		return ret, nil
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", cloud))
	}
	ret, err = transformer.CloudToHub(input)
	return ret, err
}
func fromHubFormat(hub clusterinfo.ClusterInfo, toCloud string, outputScope string) (c clusterinfo.ClusterInfo, err error) {
	if hub.Cloud != clusterinfo.HUB {
		return clusterinfo.ClusterInfo{}, errors.New(fmt.Sprintf("Wrong Cloud %s", hub.Cloud))
	}
	var transformer clusteraccess.Transformer
	err = nil
	var ret clusterinfo.ClusterInfo
	switch toCloud { //  We do not expect more than  these clouds so not splitting out dynamically loaded adapters
	case clusterinfo.GCP:
		transformer = &transformgke.GkeTransformer{}
	case clusterinfo.AZURE:
		transformer = &transformaks.AksTransformer{}
	case clusterinfo.AWS:
		return c, errors.New(fmt.Sprintf("Unsupported %s", toCloud))
	case clusterinfo.HUB:
		log.Print("to Hub , no changes")
		ret = hub
		return ret, nil
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", toCloud))
	}
	ret, err = transformer.HubToCloud(hub, outputScope)
	return ret, err
}
