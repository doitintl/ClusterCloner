package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"github.com/stretchr/testify/assert"
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
