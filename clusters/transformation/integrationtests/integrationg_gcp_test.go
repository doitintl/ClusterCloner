package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestCreateGCPClusterFromFile ...
func TestCreateGCPClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/gke_clusters.json"
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster
	out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, clusters.GCP, scopeForTest, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
		t.Errorf("%s  does not have %s as prefix", out[0].Name, clusterFromFile.Name)
	}

	// assertNumberClustersByLabel the created files, by label
	ca := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	assertNumberClustersByLabel(t, clusterFromFile, 1)

	//Delete it
	err = ca.Delete(out[0])
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClustersByLabel(t, clusterFromFile, 0)

}
