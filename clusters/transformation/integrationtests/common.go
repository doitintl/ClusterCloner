package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var scopeForTest = "joshua-playground" // TODO parametrize

func getCreatedClustersByLabel(t *testing.T, searchTemplate *clusters.ClusterInfo, expectedCount int) *clusters.ClusterInfo {
	ca := clusteraccess.GetClusterAccess(searchTemplate.Cloud)
	listed, err := ca.List(searchTemplate.Scope, searchTemplate.Location, searchTemplate.Labels)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedCount, len(listed), listed)
	if len(listed) == 0 {
		return nil
	}
	return listed[0]
}

// execTestClusterFromFile ...
func execTestClusterFromFile(t *testing.T, inputFile string) {
	clustersFromFile := cleanAndCreateClusterFromFile(t, inputFile)
	deleteAllMatchingByLabel(t, clustersFromFile)

}

func cleanAndCreateClusterFromFile(t *testing.T, inputFile string) []*clusters.ClusterInfo {
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	deleteAllMatchingByLabel(t, clustersFromFile)

	for _, clusterFromFile := range clustersFromFile {
		inCloud := clusterFromFile.Cloud
		outCloud := inCloud
		out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, outCloud, scopeForTest, true, true)
		if err != nil {
			t.Fatal(err)
		}
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

		c.Scope = scopeForTest
		ca := clusteraccess.GetClusterAccess(c.Cloud)
		listed, err := ca.List(c.Scope, c.Location, c.Labels)
		if err != nil {
			t.Fatal(err)
		}
		if len(listed) > 0 {
			for _, deleteThis := range listed {
				err = ca.Delete(deleteThis)
			}
			if err != nil {
				t.Fatal(err)
			}
		}
		_ = getCreatedClustersByLabel(t, c, 0)
	}
}
