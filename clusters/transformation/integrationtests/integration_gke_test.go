package integrationtests

import (
	"testing"
)

// TestCreateGCPClusterFromFile ...
func TestCreateGCPClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	execTestClusterFromFile(t, "test-data/gke_clusters.json")
}
