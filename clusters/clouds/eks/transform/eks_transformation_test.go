package transform

import (
	"clustercloner/clusters/clouds/gke/transform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformLocToHub(t *testing.T) {
	loc := "us-east-2"
	locationMap, err := LocationsCloudToHub()
	assert.Nil(t, err)
	hub := locationMap[loc]
	assert.Equal(t, "us-central1", hub)
	gcpLoc, err := transform.LocationsCloudToHub()
	assert.Nil(t, err)
	for _, gcp := range locationMap {
		_, ok := gcpLoc[gcp]
		if !ok {
			t.Fatal(gcp)
		}
	}
}
