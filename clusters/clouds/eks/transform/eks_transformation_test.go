package transform

import (
	"clustercloner/clusters/clouds/gke/transform"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformLocToHub(t *testing.T) {
	loc := "us-east-2"
	locationMap, err := getAwsToHubLocations()
	if err != nil {
		t.Fatal(err)
	}
	hub := locationMap[loc]
	assert.Equal(t, "us-central1", hub)
	gcpLoc, err := transform.GetGcpLocations()
	if err != nil {
		t.Fatal(err)
	}
	for _, gcp := range locationMap {
		if !util.ContainsStr(gcpLoc, gcp) {
			t.Error(gcp)
		}
	}
}
