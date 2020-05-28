package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var scopeForTest = "joshua-playground"

func assertNumberClustersByLabel(t *testing.T, ci *clusters.ClusterInfo, expected int) {
	ca := clusteraccess.GetClusterAccess(ci.Cloud)
	listed, err := ca.List(ci.Scope, ci.Location, ci.Labels)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, len(listed), listed)
}

// execTestClusterFromFile ...
func execTestClusterFromFile(t *testing.T, inputFile string) {
	clustersFromFile, err := clusters.LoadFromFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	clusterFromFile := clustersFromFile[0]
	assert.Equal(t, 1, len(clustersFromFile), "we work with a single cluster in this test")

	//create Cluster
	inCloud := clusterFromFile.Cloud
	outCloud := inCloud
	out, err := transformation.Clone(inputFile, "", "", "", clusterFromFile.Labels, outCloud, scopeForTest, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out[0].Name, clusterFromFile.Name) {
		t.Fatalf("%s does not have %s as prefix", out[0].Name, clusterFromFile.Name)
	}

	assertNumberClustersByLabel(t, clusterFromFile, 1)
	ca := clusteraccess.GetClusterAccess(clusterFromFile.Cloud)
	err = ca.Delete(out[0])
	if err != nil {
		t.Fatal(err)
	}
	assertNumberClustersByLabel(t, clusterFromFile, 0)
}
