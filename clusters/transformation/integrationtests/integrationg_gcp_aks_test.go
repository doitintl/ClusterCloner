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
func TestCreateGCPClusterFromFileThenCloneToAKS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/gke_clusters.json"
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster from file
	createdGCPClusters, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, clusters.GCP, scopeForTest, true, true)
	if err != nil {
		t.Fatal(err)
	}
	createdGCPCluster := createdGCPClusters[0]
	if !strings.HasPrefix(createdGCPCluster.Name, clusterFromFile.Name) {
		t.Errorf("%s does not have %s as prefix", createdGCPCluster.Name, clusterFromFile.Name)
	}

	// assertNumberClustersByLabel the created files, by label
	gkeAccess := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	assertNumberClustersByLabel(t, clusterFromFile, 1)
	createdAKSClusters, err := transformation.Clone("",
		createdGCPCluster.Cloud,
		createdGCPCluster.Scope,
		createdGCPCluster.Location,
		createdGCPCluster.Labels,
		clusters.Azure,
		scopeForTest,
		true,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
	createdAKSCluster := createdAKSClusters[0]
	assertNumberClustersByLabel(t, createdAKSCluster, 1)

	//Delete both
	err = gkeAccess.Delete(createdGCPCluster)
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClustersByLabel(t, createdGCPCluster, 0)

	var aksAccess = clusteraccess.GetClusterAccess(createdAKSCluster.Cloud)
	err = aksAccess.Delete(createdAKSCluster)
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClustersByLabel(t, createdAKSCluster, 0)

}
