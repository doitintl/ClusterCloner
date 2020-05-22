package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"github.com/stretchr/testify/assert"
	"testing"
)

var scopeForTest = "joshua-playground" //TODO take this from config here and in the JSONs used for testing

func assertNumberClustersByLabel(t *testing.T, ci *clusters.ClusterInfo, expected int) {
	ca := clusteraccess.GetClusterAccess(ci.Cloud)
	listed, err := ca.List(ci.Scope, ci.Location, ci.Labels)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, len(listed), listed)
}
