package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/gke/transform"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransformAzureToHub(t *testing.T) {
	ci := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.Azure,
		Location:    "westus2",
		Scope:       "samplescope",
		K8sVersion:  "1.14.0",
		Labels:      map[string]string{"a": "aa", "b": "bb"},
		GeneratedBy: clusters.Mock}
	tr := AKSTransformer{}
	std, err := tr.CloudToHub(ci)
	assert.Nil(t, err)
	if std.Location != "us-west1" {
		t.Fatal(std.Location)
	}
	if std.Cloud != clusters.Hub {
		t.Fatalf("not the standard cloud %s", std.Cloud)
	}
}
func TestTransformAzureToHubBadLoc(t *testing.T) {
	ci := &clusters.ClusterInfo{Name: "c",
		Cloud: clusters.Azure, Location: "westus1",
		Scope:       "sampelscope",
		K8sVersion:  "1.15.0",
		GeneratedBy: clusters.Mock,
		Labels:      map[string]string{"a": "aa", "b": "bb"},
	}
	tr := AKSTransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestTransformHubToAzure(t *testing.T) {
	ci := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.Hub,
		Location:    "us-central1",
		Scope:       "",
		K8sVersion:  "1.14.6",
		Labels:      map[string]string{"a": "aa", "b": "bb"},
		GeneratedBy: clusters.Mock,
	}
	tr := AKSTransformer{}
	az, err := tr.HubToCloud(ci, "")
	assert.Nil(t, err)
	if az.Location != "centralus" {
		t.Fatal(az.Location)
	}
	if az.Cloud != clusters.Azure {
		t.Fatalf("Not the expected cloud %s", az.Cloud)
	}
	if !strings.HasPrefix(az.K8sVersion, "1.14") {
		t.Fatalf("Bad K8s Version for Azure based on input: %s", az.K8sVersion)
	}
}

func TestTransformLocToHub(t *testing.T) {
	loc := "eastus"
	locationMap, err := locationsCloudToHub()
	assert.Nil(t, err)
	hub, ok := locationMap.Get(loc)
	assert.True(t, ok)
	assert.Equal(t, "us-east4", hub)
	_, err = transform.LocationsCloudToHub()
	assert.Nil(t, err)
}
