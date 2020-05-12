package transform

import (
	"clustercloner/clusters"
	"testing"
)

func TestTransformAzureToHub(t *testing.T) {
	ci := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.Azure,
		Location:    "westus2",
		Scope:       "samplescope",
		K8sVersion:  "1.14.0",
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
		GeneratedBy: clusters.Mock}
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
