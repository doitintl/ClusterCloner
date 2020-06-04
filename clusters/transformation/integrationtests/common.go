package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"clustercloner/clusters/transformation/util"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func getCreatedClustersByLabel(t *testing.T, searchTemplate *clusters.ClusterInfo, expectedCount int) *clusters.ClusterInfo {
	ca := clusteraccess.GetClusterAccess(searchTemplate.Cloud)
	listed, err := ca.List(searchTemplate.Scope, searchTemplate.Location, searchTemplate.Labels)
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, len(listed), listed)
	if len(listed) == 0 {
		return nil
	}
	return listed[0]
}

func cleanCreateDeleteCluster(t *testing.T, inputFile string, alsoDescribe bool) {
	createdClusters := cleanAndCreateCluster(t, inputFile)
	defer deleteAllMatchingByLabel(t, createdClusters)
	if alsoDescribe {
		for _, cl := range createdClusters {
			ca := clusteraccess.GetClusterAccess(cl.Cloud)
			searchTemplate := clusters.ClusterInfo{Cloud: cl.Cloud, Scope: cl.Scope, Location: cl.Location, Name: cl.Name, GeneratedBy: clusters.SearchTemplate}
			description, err := ca.Describe(&searchTemplate)
			if err != nil && strings.Contains(err.Error(), "not found") {
				assert.Fail(t, "cluster not found "+err.Error())
			}
			assert.Nil(t, err)
			assert.Equal(t, cl, description)
		}
	}

}

func cleanAndCreateCluster(t *testing.T, inputFile string) []*clusters.ClusterInfo {
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	assert.Nil(t, err)

	deleteAllMatchingByLabel(t, clustersFromFile)

	for _, clusterFromFile := range clustersFromFile {
		inCloud := clusterFromFile.Cloud
		outCloud := inCloud
		out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, outCloud, scopeForIntegrationTest(), true, true)
		assert.Nil(t, err)
		if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
			t.Fatalf("%s does not have %s as prefix", out[0].Name, clusterFromFile.Name)
		}
		createdClus := getCreatedClustersByLabel(t, clusterFromFile, 1)
		assert.Equal(t, len(clusterFromFile.NodePools), len(createdClus.NodePools))
		for _, createdNP := range createdClus.NodePools {
			inputNP := clusterFromFile.NodePools[0]
			assert.Equal(t, inputNP.Name, createdNP.Name)
			assert.Equal(t, inputNP.DiskSizeGB, createdNP.DiskSizeGB)
			assert.Equal(t, inputNP.NodeCount, createdNP.NodeCount)
		}
	}
	return clustersFromFile
}

func deleteAllMatchingByLabel(t *testing.T, clusters []*clusters.ClusterInfo) {
	for _, c := range clusters {

		c.Scope = scopeForIntegrationTest()
		ca := clusteraccess.GetClusterAccess(c.Cloud)
		listed, err := ca.List(c.Scope, c.Location, c.Labels)
		assert.Nil(t, err)
		if len(listed) > 0 {
			for _, deleteThis := range listed {
				err = ca.Delete(deleteThis)
			}
			assert.Nil(t, err)
		}
		_ = getCreatedClustersByLabel(t, c, 0)
	}
}

func scopeForIntegrationTest() string {
	key := "AZURE_BASE_GROUP_NAME"
	val := os.Getenv(key)
	if val == "" {
		panic("cannot run integration tests; need to define \"scope\" (Azure Group and Google Project name) in environment variable " + key + " (in the .env file)")
	}
	return val
}
func runClusterCloning(t *testing.T, file string, outputCloud string) {
	//delete any stray input clusters, then create the input clusters
	inputClusters := cleanAndCreateCluster(t, file)
	//delete any stray potential target clusters
	for _, inCluster := range inputClusters {
		potentialTargetCluster := util.CopyClusterInfo(inCluster)
		potentialTargetCluster.Cloud = outputCloud
		potentialTargetClusters := []*clusters.ClusterInfo{&potentialTargetCluster}
		deleteAllMatchingByLabel(t, potentialTargetClusters)
	}
	//create target clustes
	created := make([]*clusters.ClusterInfo, 0)
	for _, inputCluster := range inputClusters {
		cloneOutput, err := transformation.Clone("",
			inputCluster.Cloud,
			inputCluster.Scope,
			inputCluster.Location,
			inputCluster.Labels,
			outputCloud,
			scopeForIntegrationTest(),
			true,
			true,
		)

		created = append(created, cloneOutput...)
		assert.Nil(t, err)
	}
	deleteAllMatchingByLabel(t, inputClusters)
	deleteAllMatchingByLabel(t, created)
}
