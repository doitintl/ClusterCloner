package transformation

import (
	"clustercloner/clusters"
	transformaks "clustercloner/clusters/clouds/aks/transform"
	transformgke "clustercloner/clusters/clouds/gke/transform"
	"clustercloner/clusters/clusteraccess"
	clusterutil "clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	"log"
)

// Transformer ...
type Transformer interface {
	CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error)
	//	HubToCloud///
	HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error)
	// LocationHubToCloud ...
	LocationHubToCloud(loc string) (string, error)
	// LocationCloudToHub ...
	LocationCloudToHub(loc string) (string, error)
}

func getTransformer(cloud string) Transformer {
	var transformer Transformer
	switch cloud {
	case clusters.GCP:
		transformer = &transformgke.GKETransformer{}
	case clusters.AZURE:
		transformer = &transformaks.AKSTransformer{}
	case clusters.HUB:
		transformer = &IdentityTransformer{clusters.HUB}
	default:
		transformer = nil
		log.Printf("Unknown %s", cloud)
	}
	return transformer
}

// Clone ...
func Clone(cliCtx *cli.Context) ([]*clusters.ClusterInfo, error) {
	inputCloud := cliCtx.String("inputcloud")
	outputCloud := cliCtx.String("outputcloud")
	inputLocation := cliCtx.String("inputlocation")
	inputScope := cliCtx.String("inputscope")
	outputScope := cliCtx.String("outputscope")
	randSfx := cliCtx.Bool("randomsuffix")
	create := cliCtx.Bool("create")
	sCreate := ""
	if !create {
		sCreate = "not "
	}
	log.Printf("Will %screate target clusters", sCreate)
	if inputCloud == "" || outputCloud == "" || inputLocation == "" || inputScope == "" || outputScope == "" {
		log.Fatal("Missing flags")
	}
	return clone(inputCloud, outputCloud, inputLocation, inputScope, outputScope, create, randSfx)

}

func clone(inputCloud string, outputCloud string, inputLocation string, inputScope string, outputScope string, create bool, randSfx bool) ([]*clusters.ClusterInfo, error) {
	clusterAccessor := clusteraccess.GetClusterAccess(inputCloud)
	if clusterAccessor == nil {
		return nil, errors.New("cannot get accessor for " + inputCloud)
	}
	inputClusterInfos, err := clusterAccessor.ListClusters(inputScope, inputLocation)
	if err != nil {
		return nil, err
	}
	transformationOutput := make([]*clusters.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := transformCloudToCloud(inputClusterInfo, outputCloud, outputScope, randSfx)
		transformationOutput = append(transformationOutput, outputClusterInfo)
		if err != nil {
			log.Printf("Error processing %v: %v", inputClusterInfo, err)
		}
	}

	if create {
		createdClusters, createdIndexes, _ := createClusters(transformationOutput)
		//replaced each ClusterInfo that was created; left the ones that were not
		var transformedClustersCreatedOrNot = make([]*clusters.ClusterInfo, len(transformationOutput))
		if len(createdClusters) != len(transformationOutput) {
			panic(fmt.Sprintf("%d!=%d", len(createdClusters), len(transformationOutput)))
		}
		for i := 0; i < len(transformationOutput); i++ {
			if clusterutil.ContainsInt(createdIndexes, i) {
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
func createClusters( /*immutable*/ createThese []*clusters.ClusterInfo) (createdClusters []*clusters.ClusterInfo, createdIndexes []int, err error) {
	createdIndexes = make([]int, 0)
	createdClusters = make([]*clusters.ClusterInfo, len(createThese))
	log.Println("Creating", len(createThese), "target clusters")
	for idx, createThis := range createThese {
		created, err := createCluster(createThis)
		if err != nil {
			log.Printf("Error creating %v: %v", createThis, err)
		} else {
			createdClusters[idx] = created //TODO read the clusters back from the cloud, where feasible (but async creation may prevent that or require delay)
			createdIndexes = append(createdIndexes, idx)
		}
	}
	return createdClusters, createdIndexes, nil
}

// createCluster ...
func createCluster(createThis *clusters.ClusterInfo) (createdClusterInfo *clusters.ClusterInfo, err error) {
	var ca = clusteraccess.GetClusterAccess(createThis.Cloud)
	if ca == nil {
		return nil, errors.New("cannot creeate ClusterAccess")
	}
	created, err := ca.CreateCluster(createThis)
	if err != nil {
		log.Println("Error creating cluster", err)
	}
	return created, err
}

func transformCloudToCloud(in *clusters.ClusterInfo, toCloud, outputScope string, randSfx bool) (c *clusters.ClusterInfo, err error) {
	if in.Cloud == toCloud {
		t := IdentityTransformer{toCloud}
		out, err := t.HubToCloud(in, outputScope)
		if err != nil || out == nil {
			return nil, errors.Wrap(err, "Error in transforming to same cloud")
		}
		var sfx string
		if randSfx {
			sfx = clusterutil.RandomAlphaNumSequence(5, false, true, true)
		} else {
			sfx = "copy"
		}
		out.Name = out.Name + "-" + sfx
		return out, nil
	}
	hub, err1 := toHubFormat(in)
	if err1 != nil || hub == nil {
		return nil, errors.Wrap(err1, "Error in transforming toHubFormat")
	}
	out, err2 := fromHubFormat(hub, toCloud, outputScope, randSfx)
	if err2 != nil {
		return nil, errors.Wrap(err2, "cannot convert from Hub format")
	}
	out.GeneratedBy = clusters.TRANSFORMATION
	return out, nil

}

func toHubFormat(input *clusters.ClusterInfo) (ret *clusters.ClusterInfo, err error) {
	cloud := input.Cloud
	transformer := getTransformer(cloud)
	if transformer == nil {
		return nil, errors.New("cannot transform")
	}
	ret, err = transformer.CloudToHub(input)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert CloudToHub")
	}

	return ret, err
}

func fromHubFormat(hub *clusters.ClusterInfo, toCloud string, outputScope string, randSuffix bool) (ret *clusters.ClusterInfo, err error) {
	if hub.Cloud != clusters.HUB {
		return nil, errors.New(fmt.Sprintf("Wrong Cloud %s", hub.Cloud))
	}

	var transformer = getTransformer(toCloud)
	ret, err = transformer.HubToCloud(hub, outputScope)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert HubToCloud")
	}
	if randSuffix {
		ret.Name = ret.Name + "-" + clusterutil.RandomAlphaNumSequence(5, false, true, true)
	}

	return ret, err
}
