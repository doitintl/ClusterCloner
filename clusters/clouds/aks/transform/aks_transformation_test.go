package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/gke/transform"
	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Error(err)
	}
	if std.Location != "us-west1" {
		t.Error(std.Location)
	}
	if std.Cloud != clusters.Hub {
		t.Errorf("not the standard cloud %s", std.Cloud)
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
		t.Error("expect error")
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
	if err != nil {
		t.Error(err)
	}
	if az.Location != "centralus" {
		t.Error(az.Location)
	}
	if az.Cloud != clusters.Azure {
		t.Errorf("Not the expected cloud %s", az.Cloud)
	}
	if az.K8sVersion != "1.14.7" {
		t.Errorf("Bad K8s Version for Azure based on input: %s", az.K8sVersion)
	}
}

func TestTransformLocToHub(t *testing.T) {
	loc := "eastus"
	locationMap, err := LocationsCloudToHub()
	if err != nil {
		t.Fatal(err)
	}
	hub := locationMap[loc]
	assert.Equal(t, "us-east4", hub)
	gcpLoc, err := transform.LocationsCloudToHub()
	if err != nil {
		t.Fatal(err)
	}
	for _, gcp := range locationMap {
		if _, ok := gcpLoc[gcp]; !ok {
			t.Error(gcp)
		}
	}
}
