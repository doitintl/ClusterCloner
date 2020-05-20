package transformation

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateGCPClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/gke_clusters.json"
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster
	out, err := Clone(inputFile, "", "", "", clusterFromFile.Labels, clusters.GCP, "joshua-playground", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
		t.Errorf("%s != %s", out[0].Name, clusterFromFile.Name)
	}

	// assertNumberClusters the created files, by label
	ca := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	assertNumberClusters(t, clusterFromFile, 1)

	//Delete it
	err = ca.Delete(out[0])
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClusters(t, clusterFromFile, 0)

}

func assertNumberClusters(t *testing.T, ci *clusters.ClusterInfo, expected int) {
	ca := clusteraccess.GetClusterAccess(ci.Cloud)
	// assertNumberClusters the created files, after deletion. There should be none
	listedAfterDel, err := ca.List(ci.Scope, ci.Location, ci.Labels)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, len(listedAfterDel), listedAfterDel)
}
