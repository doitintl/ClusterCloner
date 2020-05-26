package integrationtests

import (
	"testing"
)

// TestCreateGCPClusterFromFile ...
func TestCreateGCPClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/gke_clusters.json"
	CreateClusterFromFile(t, inputFile)
}
