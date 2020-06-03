package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/transformation"
	"testing"
)

// TestCreateGCPClusterFromFile ...
func TestCreateGCPClusterFromFileThenCloneToAKS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	created := make([]*clusters.ClusterInfo, 0)
	createdGCPClusters := cleanAndCreateClusterFromFile(t, "test-data/gke_clusters.json")
	for _, createdGCPCluster := range createdGCPClusters {

		cloneOutput, err := transformation.Clone("",
			createdGCPCluster.Cloud,
			createdGCPCluster.Scope,
			createdGCPCluster.Location,
			createdGCPCluster.Labels,
			clusters.Azure,
			scopeForTest,
			true,
			true,
		)
		created = append(created, cloneOutput...)
		if err != nil {
			t.Fatal(err)
		}

	}
	deleteAllMatchingByLabel(t, createdGCPClusters)
	deleteAllMatchingByLabel(t, created)
}
