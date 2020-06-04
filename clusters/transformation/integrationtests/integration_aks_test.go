package integrationtests

import (
	"testing"
)

// TestCreateAzureClusterFromFile ...
func TestCreateAzureClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cleanCreateDeleteCluster(t, "test-data/aks_clusters.json", true)

}
