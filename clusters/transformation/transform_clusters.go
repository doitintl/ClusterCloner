package transformation

import (
	"clustercloner/clusters"
	transformaks "clustercloner/clusters/clouds/aks/transform"
	transformeks "clustercloner/clusters/clouds/eks/transform"
	transformgke "clustercloner/clusters/clouds/gke/transform"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation/util"
	clusterutil "clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
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

// getTransformer, from or to hub to the cloud specified by the parameter.
func getTransformer(cloud string) (Transformer, error) {
	switch cloud {
	case clusters.GCP:
		return &transformgke.GKETransformer{}, nil
	case clusters.Azure:
		return &transformaks.AKSTransformer{}, nil
	case clusters.AWS:
		return &transformeks.EKSTransformer{}, nil
	default:
		return nil, errors.New("Unknown cloud \"" + cloud + "\"")
	}
}

// getTransformer, from this cloud to this same cloud, without going through hub
func getSameCloudTransformer(cloud string) Transformer {
	var transformer Transformer
	switch cloud {
	case clusters.GCP:
		transformer = &transformgke.GKEToGKETransformer{}
	case clusters.Azure:
		transformer = &util.IdentityTransformer{TargetCloud: clusters.Azure}
	case clusters.AWS:
		transformer = &util.IdentityTransformer{TargetCloud: clusters.AWS}
	case clusters.Hub:
		transformer = &util.IdentityTransformer{TargetCloud: clusters.Hub}
	default:
		transformer = nil
		panic(fmt.Sprintf("cannot get transformer for unknown cloud %s", cloud))
	}
	return transformer
}

// CloneFromCli ...
func CloneFromCli(cliCtx *cli.Context) ([]*clusters.ClusterInfo, error) {
	inputFile, inputCloud, outputCloud, inputLocation, inputScope, outputScope, shouldCreate, randSfx, labelFilter, err := parseCLIParams(cliCtx)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse CLI params")
	}
	return Clone(inputFile, inputCloud, inputScope, inputLocation, labelFilter, outputCloud, outputScope, randSfx, shouldCreate)

}

// Clone ...
func Clone(inputFile string, inputCloud string, inputScope string, inputLocation string, labelFilter map[string]string, outputCloud string, outputScope string, randSfx bool, shouldCreate bool) ([]*clusters.ClusterInfo, error) {
	if labelFilter == nil { // We usually enforce non-nil labelFilter to make sure we copy it. Here, nil is acceptable
		labelFilter = make(map[string]string)
	}
	inputClusterInfos, err := getInputClusters(inputFile, inputCloud, inputScope, inputLocation, labelFilter)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get input clusters")
	}
	transformationOutput, err := transform(inputClusterInfos, outputCloud, outputScope, randSfx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot transform clusters")
	}
	if len(transformationOutput) != len(inputClusterInfos) {
		panic(fmt.Sprintf("%d!=%d", len(transformationOutput), len(inputClusterInfos)))
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

func parseCLIParams(cliCtx *cli.Context) (inputFile string, inputCloud string, outputCloud string, inputLocation string, inputScope string, outputScope string, shouldCreate, randomSuffix bool, labelFilter map[string]string, err error) {
	inputFile = cliCtx.String("inputfile")
	inputCloud = cliCtx.String("inputcloud")
	outputCloud = cliCtx.String("outputcloud")
	inputLocation = cliCtx.String("inputlocation")
	inputScope = cliCtx.String("inputscope")
	outputScope = cliCtx.String("outputscope")
	randomSuffix = cliCtx.Bool("randomsuffix")

	shouldCreate = cliCtx.Bool("nodryrun")
	if shouldCreate {
		log.Printf("No dry run; will create target clusters")
	} else {
		log.Printf("Dry run; will not create target clusters")
	}
	labelFilterS := cliCtx.String("labelfilter")
	labelFilter = clusterutil.CommaSeparatedKeyValPairsToMap(labelFilterS)
	var errS string
	if inputFile == "" {
		if inputCloud == "" {
			errS += "; input cloud missing"
		}
		if inputScope == "" && inputCloud != clusters.AWS {
			errS += "; input scope missing"
		}
		if inputLocation == "" {
			errS += "; input location missing"
		}
	} else {
		if inputCloud != "" || inputScope != "" || inputLocation != "" {
			errS += "; if input file is provided, do not provide input cloud name, scope, or location"
		}
	}
	if (inputScope != "" && inputCloud != clusters.AWS) || (outputScope != "" && outputCloud != clusters.AWS) {
		errS += "; do not specify scope for AWS"
	}

	if outputCloud == "" {
		errS += "; output cloud missing"
	}
	if outputScope == "" && outputCloud != clusters.AWS {
		errS += "; output scope"
	}

	if errS != "" {
		errS = strings.TrimPrefix(errS, "; ")
		err = errors.New(errS)
	}

	return inputFile, inputCloud, outputCloud, inputLocation, inputScope, outputScope, shouldCreate, randomSuffix, labelFilter, err
}

func transform(inputClusterInfos []*clusters.ClusterInfo, outputCloud string, outputScope string, randSfx bool) ([]*clusters.ClusterInfo, error) {
	transformationOutput := make([]*clusters.ClusterInfo, 0)
	for _, inputClusterInfo := range inputClusterInfos {
		outputClusterInfo, err := transformCloudToCloud(inputClusterInfo, outputCloud, outputScope, randSfx)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error processing %v: %v", inputClusterInfo, err))
		}
		if outputClusterInfo == nil {
			panic("should not be nil")
		}
		validateSourceCluster(outputClusterInfo, clusters.Transformation)
		transformationOutput = append(transformationOutput, outputClusterInfo)

	}
	return transformationOutput, nil
}

func getInputClusters(inputFile string, inputCloud string, inputScope string, inputLocation string, labelFilter map[string]string) (inputClusterInfos []*clusters.ClusterInfo, err error) {
	if inputFile != "" {
		inputClusterInfos, err = clusters.LoadFromFile(inputFile)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read input file")
		}
		log.Printf("Loaded %d clusters from file\n", len(inputClusterInfos))
	} else {
		clusterAccessor := clusteraccess.GetClusterAccess(inputCloud)
		if clusterAccessor == nil {
			return nil, errors.New("cannot get accessor for " + inputCloud)
		}
		inputClusterInfos, err = clusterAccessor.List(inputScope, inputLocation, labelFilter)
		if err != nil {
			return nil, errors.Wrap(err, "error listing clusters")
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

// 'createdClusters' return param may be partly populated as some clusters have been successfully populated and some have not
func createClusters0( /*immutable*/ createThese []*clusters.ClusterInfo) (createdClusters []*clusters.ClusterInfo, createdIndexes []int) {
	createdIndexes = make([]int, 0)
	createdClusters = make([]*clusters.ClusterInfo, len(createThese))
	log.Println("Creating", len(createThese), "clusters")
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

func createCluster(createThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error) {

	var ca = clusteraccess.GetClusterAccess(createThis.Cloud)
	if ca == nil {
		return nil, errors.New("cannot create ClusterAccess")
	}
	created, err = ca.Create(createThis)
	if err != nil {
		return nil, errors.Wrap(err, "error creating cluster "+createThis.Name)
	}
	validateSourceCluster(created, clusters.Created)

	return created, nil
}

func validateSourceCluster(ci *clusters.ClusterInfo, expectedGenByForCluster string) {
	if ci.GeneratedBy != expectedGenByForCluster {
		log.Printf("Error: Actual %s != expected %s", ci.GeneratedBy, expectedGenByForCluster)
	}
	if ci.Labels == nil {
		panic("Must initialize Labels")
	}

	var expectedGenByForSource []string = nil
	switch ci.GeneratedBy {
	case clusters.Mock:
		expectedGenByForSource = []string{""}
	case clusters.Read:
		expectedGenByForSource = []string{clusters.SearchTemplate, ""}
	case clusters.Created:
		expectedGenByForSource = []string{clusters.Transformation}
	case clusters.Transformation:
		//Source of Transformation is Transformation when we have transformed twice: to hub and from hub
		expectedGenByForSource = []string{clusters.Transformation, clusters.Read, clusters.InputFile}
	case clusters.InputFile:
		expectedGenByForSource = nil //nil means "don't check"
	case clusters.SearchTemplate:
		expectedGenByForSource = []string{""}
	default:
		panic("unknown " + ci.GeneratedBy)
	}

	if ci.SourceCluster == nil {
		return
	}
	actual := ci.SourceCluster.GeneratedBy

	actualIsExpected := false
	if expectedGenByForSource == nil { //nil means "don't check"
		actualIsExpected = true
	} else {
		actualIsExpected = clusterutil.ContainsStr(expectedGenByForSource, actual)
	}
	if !actualIsExpected {
		log.Printf("unexpected GeneratedBy for SourceCluster: \"%s\" is not one of \"%s\"\n%s", actual, expectedGenByForSource, clusterutil.ToJSON(ci))
	} else {
		validateSourceCluster(ci.SourceCluster, ci.SourceCluster.GeneratedBy)
	}

}

func transformCloudToCloud(in *clusters.ClusterInfo, toCloud, outputScope string, randSfx bool) (c *clusters.ClusterInfo, err error) {
	var out *clusters.ClusterInfo
	if in.Cloud == toCloud { //don't use hub
		t := getSameCloudTransformer(toCloud)
		out, err = t.HubToCloud(in, outputScope)
		if err != nil {
			return nil, errors.Wrap(err, "Error in transforming to same cloud")
		}
		var sfx string
		if randSfx {
			sfx = clusterutil.RandomWord()
		} else {
			sfx = "copy"
		}
		out.Name = out.Name + "-" + sfx
	} else { //two different clouds;so we use use hub
		hub, err := toHubFormat(in)
		if err != nil || hub == nil {
			return nil, errors.Wrap(err, "error in transforming toHubFormat")
		}
		out, err = fromHubFormat(hub, toCloud, outputScope, randSfx)
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert from Hub format")
		}
	}

	out.Labels["clustercloner-source-cloud"] = clusterutil.ToLowerCaseAlphaNumDashAndUnderscore(in.Cloud)
	out.GeneratedBy = clusters.Transformation
	return out, nil

}

func toHubFormat(input *clusters.ClusterInfo) (ret *clusters.ClusterInfo, err error) {
	cloud := input.Cloud
	transformer, err := getTransformer(cloud)
	if err != nil {
		return nil, errors.Wrap(err, "cannot getTransformer to convert toHubFormat")
	}

	ret, err = transformer.CloudToHub(input)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert CloudToHub")
	}

	return ret, nil
}

func fromHubFormat(hub *clusters.ClusterInfo, toCloud string, outputScope string, randSuffix bool) (ret *clusters.ClusterInfo, err error) {
	if hub.Cloud != clusters.Hub {
		return nil, errors.Errorf("wrong Cloud %s", hub.Cloud)
	}

	transformer, err := getTransformer(toCloud)
	if err != nil {
		return nil, errors.Wrap(err, "cannot getTransformer to convert fromHubFormat")
	}
	ret, err = transformer.HubToCloud(hub, outputScope)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert HubToCloud")
	}
	if randSuffix {
		ret.Name = ret.Name + "-" + clusterutil.RandomWord()
	}

	return ret, nil
}
