package integrationtests

import (
	"testing"
)

// TestCreateAzureClusterFromFile ...
func TestCreateAzureClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	execTestClusterFromFile(t, "test-data/aks_clusters.json")

}
