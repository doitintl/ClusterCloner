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
	case clusters.Azure:
		transformer = &transformaks.AKSTransformer{}
	case clusters.Hub:
		transformer = &IdentityTransformer{clusters.Hub}
	default:
		transformer = nil
		log.Printf("Unknown %s", cloud)
	}
	return transformer
}

// Clone ...
func Clone(cliCtx *cli.Context) ([]*clusters.ClusterInfo, error) {
	inputFile := cliCtx.String("inputfile")
	inputCloud := cliCtx.String("inputcloud")
	outputCloud := cliCtx.String("outputcloud")
	inputLocation := cliCtx.String("inputlocation")
	inputScope := cliCtx.String("inputscope")
	outputScope := cliCtx.String("outputscope")
	randSfx := cliCtx.Bool("randomsuffix")
	if (inputFile != "") != (inputCloud != "" || inputLocation != "" || inputScope != "") {
		return nil, errors.New("Must read either from file or from cloud, not both")
	}
	// We could validate that if we read from cloud, all relevant input CLI params must be there,
	// and likewise if create is true, then all output CLI params are there.

	create := cliCtx.Bool("create")
	if create {
		log.Printf("Will create target clusters")
	} else {
		log.Printf("Dry run; will not create target clusters")
	}

	return clone(inputFile, inputCloud, outputCloud, inputLocation, inputScope, outputScope, create, randSfx)

}

func clone(inputFile string,
	inputCloud string,
	outputCloud string,
	inputLocation string,
	inputScope string,
	outputScope string,
	shouldCreate bool,
	randSfx bool) (transformedOrCreated []*clusters.ClusterInfo, err error) {

	inputClusterInfos, err := getInputClusters(inputFile, inputCloud, err, inputScope, inputLocation)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get input clusters")
	}
	transformationOutput, err := transform(inputClusterInfos, outputCloud, outputScope, randSfx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot transform clusters")
	}

	if !shouldCreate {
		log.Println("Dry run, not creating", len(transformationOutput), "target clusters")
		return transformationOutput, nil
	}

	transformedClustersCreatedOrNot, err := createClusters(transformationOutput)
	if err != nil {
		return nil, errors.Wrap(err, "cannot shouldCreate")
	}
	return transformedClustersCreatedOrNot, nil

}

func transform(inputClusterInfos []*clusters.ClusterInfo, outputCloud string, outputScope string, randSfx bool) ([]*clusters.ClusterInfo, error) {
	transformationOutput := make([]*clusters.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		assertSourceCluster(inputClusterInfo, clusters.Read)
		outputClusterInfo, err := transformCloudToCloud(inputClusterInfo, outputCloud, outputScope, randSfx)
		assertSourceCluster(outputClusterInfo, clusters.Transformation)
		transformationOutput = append(transformationOutput, outputClusterInfo)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error processing %v: %v", inputClusterInfo, err))
		}
	}
	return transformationOutput, nil
}

func getInputClusters(inputFile string, inputCloud string, err error, inputScope string, inputLocation string) ([]*clusters.ClusterInfo, error) {
	var inputClusterInfos []*clusters.ClusterInfo
	if inputFile != "" {
		panic("TODO")
	} else {

		clusterAccessor := clusteraccess.GetClusterAccess(inputCloud)
		if clusterAccessor == nil {
			return nil, errors.New("cannot get accessor for " + inputCloud)
		}
		inputClusterInfos, err = clusterAccessor.ListClusters(inputScope, inputLocation)
		if err != nil {
			return nil, err
		}
	}
	return inputClusterInfos, nil
}

func createClusters(transformationOutput []*clusters.ClusterInfo) ([]*clusters.ClusterInfo, error) {
	createdClusters, createdIndexes := createClusters0(transformationOutput)
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

// 'created' return param may be partly populated as some clusters have been successfully populated and some have not
func createClusters0( /*immutable*/ createThese []*clusters.ClusterInfo) (createdClusters []*clusters.ClusterInfo, createdIndexes []int) {
	createdIndexes = make([]int, 0)
	createdClusters = make([]*clusters.ClusterInfo, len(createThese))
	log.Println("Creating", len(createThese), "target clusters")
	for idx, createThis := range createThese {
		created, err := createCluster(createThis)
		if err != nil {
			log.Printf("error creating cluster %s: %v", createThis.Name, err)
		} else {
			createdClusters[idx] = created
			createdIndexes = append(createdIndexes, idx)
		}
	}
	return createdClusters, createdIndexes
}

func createCluster(createThis *clusters.ClusterInfo) (createdClusterInfo *clusters.ClusterInfo, err error) {
	var ca = clusteraccess.GetClusterAccess(createThis.Cloud)
	if ca == nil {
		return nil, errors.New("cannot create ClusterAccess")
	}
	created, err := ca.CreateCluster(createThis)
	if err != nil {
		log.Println("error creating cluster", err)
	} else {
		assertSourceCluster(created, clusters.Created)
	}
	return created, err
}

func assertSourceCluster(ci *clusters.ClusterInfo, expectedGenByForCluster string) {
	var expectedGenByForSource []string = nil
	if ci.GeneratedBy != expectedGenByForCluster {
		panic(fmt.Sprintf("Actual %s != expected %s", ci.GeneratedBy, expectedGenByForCluster))
	}
	switch ci.GeneratedBy {
	case clusters.Mock:
		expectedGenByForSource = []string{""}
	case clusters.Read:
		expectedGenByForSource = []string{clusters.SearchTemplate, ""}
	case clusters.Created:
		expectedGenByForSource = []string{clusters.Transformation}
	case clusters.Transformation:
		expectedGenByForSource = []string{clusters.Transformation, clusters.Read}
	case clusters.SearchTemplate:
		expectedGenByForSource = []string{""}
	default:
		panic("unknown " + ci.GeneratedBy)
	}

	var actual string
	sourceCluster := ci.SourceCluster
	if sourceCluster == nil {
		actual = ""
	} else {
		actual = sourceCluster.GeneratedBy
	}

	actualIsExpected := clusterutil.ContainsStr(expectedGenByForSource, actual)
	if !actualIsExpected {
		panic(fmt.Sprintf("unexpected GeneratedBy for SourceCluster: \"%s\"!=\"%s\"\n%s", expectedGenByForSource, actual, clusterutil.ToJSON(ci)))
	} else {
		if sourceCluster != nil {
			assertSourceCluster(sourceCluster, sourceCluster.GeneratedBy)
		}
	}

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
			sfx = clusterutil.RandomAlphaNumSequence(5, false, true, false)
		} else {
			sfx = "copy"
		}
		out.Name = out.Name + "-" + sfx
		return out, nil
	}
	hub, err1 := toHubFormat(in)
	if err1 != nil || hub == nil {
		return nil, errors.Wrap(err1, "error in transforming toHubFormat")
	}
	out, err2 := fromHubFormat(hub, toCloud, outputScope, randSfx)
	if err2 != nil {
		return nil, errors.Wrap(err2, "cannot convert from Hub format")
	}
	out.GeneratedBy = clusters.Transformation
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
	if hub.Cloud != clusters.Hub {
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
