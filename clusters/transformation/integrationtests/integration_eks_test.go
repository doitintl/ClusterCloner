package integrationtests

import (
	"testing"
)

// TestCreateAWSClusterFromFile ...
func TestCreateAWSClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/eks_clusters.json"
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	createClusterFromFile(t, inputFile)

}
