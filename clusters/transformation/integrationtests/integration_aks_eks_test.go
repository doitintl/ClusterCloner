package integrationtests

import (
	"clustercloner/clusters"
	"testing"
)

func TestCreateAzureClusterFromFileThenCloneToAWS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	file := "test-data/aks_clusters.json"
	outputCloud := clusters.AWS
	runClusterCloning(t, file, outputCloud)
}
