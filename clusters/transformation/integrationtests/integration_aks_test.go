package integrationtests

import (
	"testing"
)

// TestCreateAZClusterFromFile ...
func TestCreateAzureClusterFromFile(t *testing.T) {
	var inputFile = "test-data/aks_clusters.json"
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	CreateClusterFromFile(t, inputFile)

}
