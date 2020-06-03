package integrationtests

import (
	"testing"
)

// TestCreateAWSClusterFromFile ...
func TestCreateAWSClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	creanCreateDeleteCluster(t, "test-data/eks_clusters.json")

}
