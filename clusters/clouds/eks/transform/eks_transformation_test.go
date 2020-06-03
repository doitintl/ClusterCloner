package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformLocToHub(t *testing.T) {
	loc := "us-east-2"
	locationMap, err := LocationsCloudToHub()
	assert.Nil(t, err)
	hub, wasPresent := locationMap.Get(loc)
	assert.True(t, wasPresent)

	assert.Equal(t, "us-central1", hub)
}
