package integrationtests

import (
	"testing"
)

// TestCreateAZClusterFromFile ...
func TestCreateAzureClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/aks_clusters.json"
	execTestClusterFromFile(t, inputFile)

}
