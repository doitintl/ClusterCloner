package transformation

import (
	accessaks "clustercloner/clusters/clouds/aks/access"
	transformaks "clustercloner/clusters/clouds/aks/transform"
	//accesseks "clustercloner/clusters/clouds/eks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
	transformgke "clustercloner/clusters/clouds/gke/transform"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	"log"
)

// Transformer ...
type Transformer interface {
	CloudToHub(in *clusterinfo.ClusterInfo) (*clusterinfo.ClusterInfo, error)
	//	HubToCloud///
	HubToCloud(in *clusterinfo.ClusterInfo, outputScope string) (*clusterinfo.ClusterInfo, error)
	// LocationHubToCloud ...
	LocationHubToCloud(loc string) (string, error)
	// LocationCloudToHub ...
	LocationCloudToHub(loc string) (string, error)
}

// Clone ...
func Clone(cliCtx *cli.Context) ([]*clusterinfo.ClusterInfo, error) {
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

func clone(inputCloud string, outputCloud string, inputLocation string, inputScope string, outputScope string, create bool) ([]*clusterinfo.ClusterInfo, error) {
	var clusterAccessor clusteraccess.ClusterAccess
	switch inputCloud {
	case clusterinfo.GCP:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (project) for GCP")
		}
		clusterAccessor = accessgke.GKEClusterAccess{}
	case clusterinfo.AZURE:
		if inputScope == "" {
			return nil, errors.New("Must provide inputScope (Resource Group) for Azure")
		}
		clusterAccessor = accessaks.AKSClusterAccess{}
	default:
		return nil, errors.New(fmt.Sprintf("Invalid inputcloud \"%s\"", inputCloud))
	}
	inputClusterInfos, err := clusterAccessor.ListClusters(inputScope, inputLocation)
	if err != nil {
		return nil, err
	}
	transformationOutput := make([]*clusterinfo.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := transformCloudToCloud(inputClusterInfo, outputCloud, outputScope)
		transformationOutput = append(transformationOutput, outputClusterInfo)
		if err != nil {
			log.Printf("Error processing %v: %v", inputClusterInfo, err)
		}
	}

	if create {
		createdClusters, createdIndexes, _ := createClusters(transformationOutput)
		//replaced each ClusterInfo that was created; left the ones that were not
		var transformedClustersCreatedOrNot = make([]*clusterinfo.ClusterInfo, len(transformationOutput))
		if len(createdClusters) != len(transformationOutput) {
			panic(fmt.Sprintf("%d!=%d", len(createdClusters), len(transformationOutput)))
		}
		for i := 0; i < len(transformationOutput); i++ {
			if util.ContainsInt(createdIndexes, i) {
				transformedClustersCreatedOrNot[i] = createdClusters[i]
			} else {
				transformedClustersCreatedOrNot[i] = transformationOutput[i]
			}
		}
		return transformedClustersCreatedOrNot, nil
	}
	log.Println("Dry run, not creating", len(transformationOutput), "target clusters")
	return transformationOutput, nil

}

// 'created' return param may be partly populated as some clusters have been successfully populated and some have not
func createClusters( /*immutable*/ createThese []*clusterinfo.ClusterInfo) (createdClusters []*clusterinfo.ClusterInfo, createdIndexes []int, err error) {
	createdIndexes = make([]int, 0)
	createdClusters = make([]*clusterinfo.ClusterInfo, len(createThese))
	log.Println("Creating", len(createThese), "target clusters")
	for idx, createThis := range createThese {
		created, err := CreateCluster(createThis)
		if err != nil {
			log.Printf("Error creating %v: %v", createThis, err)
		} else {
			createdClusters[idx] = created //todo read the clusters back from the cloud, where feasible (but async creation may prevent that or require delay)
			createdIndexes = append(createdIndexes, idx)
		}
	}
	return createdClusters, createdIndexes, nil
}

// CreateCluster ...
func CreateCluster(createThis *clusterinfo.ClusterInfo) (createdClusterInfo *clusterinfo.ClusterInfo, err error) {
	var ca clusteraccess.ClusterAccess
	switch createThis.Cloud {
	case clusterinfo.AZURE:
		ca = accessaks.AKSClusterAccess{}
	case clusterinfo.AWS:
		panic("AWS not implemented")
	case clusterinfo.GCP:
		ca = accessgke.GKEClusterAccess{}
	default:
		return nil, errors.New("Unsupported Cloud for creating a cluster: " + createThis.Cloud)

	}
	created, err := ca.CreateCluster(createThis)
	if err != nil {
		log.Println("Error creating cluster", err)
	}
	return created, err
}

func transformCloudToCloud(clusterInfo *clusterinfo.ClusterInfo, toCloud string, outputScope string) (c *clusterinfo.ClusterInfo, err error) {
	if clusterInfo.Cloud == "" {
		return c, errors.New("No cloud name found in input cluster")
	}
	hub, err1 := toHubFormat(clusterInfo)
	if err1 != nil {
		return nil, errors.Wrap(err1, "Error in transforming to CloudToHub Format")
	}
	target, err2 := fromHubFormat(hub, toCloud, outputScope)
	if err2 != nil {
		return nil, err2
	}
	if clusterInfo.Cloud == toCloud {
		// Maybe self-to-self transformation shoud not go thru hub format.
		target.Name = target.Name + "-copy"
	}
	target.GeneratedBy = clusterinfo.TRANSFORMATION
	return target, nil
}
func toHubFormat(input *clusterinfo.ClusterInfo) (c *clusterinfo.ClusterInfo, err error) {
	err = nil
	var ret *clusterinfo.ClusterInfo
	var transformer Transformer
	switch cloud := input.Cloud; cloud {
	case clusterinfo.GCP:
		transformer = &transformgke.GKETransformer{}
	case clusterinfo.AZURE:
		transformer = &transformaks.AKSTransformer{}
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
func fromHubFormat(hub *clusterinfo.ClusterInfo, toCloud string, outputScope string) (c *clusterinfo.ClusterInfo, err error) {
	if hub.Cloud != clusterinfo.HUB {
		return nil, errors.New(fmt.Sprintf("Wrong Cloud %s", hub.Cloud))
	}
	var transformer Transformer
	err = nil
	var ret *clusterinfo.ClusterInfo
	switch toCloud { //  We do not expect more than  these 3 clouds so not splitting out dynamically loaded adapters
	case clusterinfo.GCP:
		transformer = &transformgke.GKETransformer{}
	case clusterinfo.AZURE:
		transformer = &transformaks.AKSTransformer{}
	case clusterinfo.AWS:
		return c, errors.New(fmt.Sprintf("Unsupported %s", toCloud))
	case clusterinfo.HUB:
		log.Print("to Hub , no changes")
		ret = hub //todo implement IdentityTransformer for this, and remove duplication from toHubFormat
		return ret, nil
	default:
		return c, errors.New(fmt.Sprintf("Unknown %s", toCloud))
	}
	ret, err = transformer.HubToCloud(hub, outputScope)
	return ret, err
}
