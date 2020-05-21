package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestCreateAZClusterFromFile ...
func TestCreateAZClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/aks_clusters.json"
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster
	out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, clusters.Azure, scopeForTest, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
		t.Errorf("%s does not have %s as prefix", out[0].Name, clusterFromFile.Name)
	}

	assertNumberClustersByLabel(t, clusterFromFile, 1)

	ca := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	err = ca.Delete(out[0])
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClustersByLabel(t, clusterFromFile, 0)

}
