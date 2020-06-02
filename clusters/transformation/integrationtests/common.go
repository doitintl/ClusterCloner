package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var scopeForTest = "joshua-playground" // TODO parametrize

func createdClusterByLabel(t *testing.T, ci *clusters.ClusterInfo, expected int) *clusters.ClusterInfo {
	ca := clusteraccess.GetClusterAccess(ci.Cloud)
	listed, err := ca.List(ci.Scope, ci.Location, ci.Labels)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, len(listed), listed)
	if len(listed) == 0 {
		return nil
	}
	return listed[0]
}

// execTestClusterFromFile ...
func execTestClusterFromFile(t *testing.T, inputFile string) {
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	for _, c := range clustersFromFile {
		c.Labels["uniqueness"] = util.RandomWord() + "-" + util.RandomWord()
	}
	if err != nil {
		t.Fatal(err)
	}
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster
	inCloud := clusterFromFile.Cloud
	outCloud := inCloud
	out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, outCloud, scopeForTest, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
		t.Fatalf("%s does not have %s as prefix", out[0].Name, clusterFromFile.Name)
	}
	createdClus := createdClusterByLabel(t, clusterFromFile, 1)
	assert.Equal(t, len(clusterFromFile.NodePools), len(createdClus.NodePools))
	for _, createdNP := range createdClus.NodePools {
		inputNP := clusterFromFile.NodePools[0]
		assert.Equal(t, inputNP.Name, createdNP.Name)
		assert.Equal(t, inputNP.DiskSizeGB, createdNP.DiskSizeGB)
		assert.Equal(t, inputNP.NodeCount, createdNP.NodeCount)
	}
	ca := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	err = ca.Delete(out[0])
	if err != nil {
		t.Fatal(err)
	}
	_ = createdClusterByLabel(t, clusterFromFile, 0)
}
