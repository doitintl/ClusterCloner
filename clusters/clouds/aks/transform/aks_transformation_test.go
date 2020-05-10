package transform

import (
	"clustercloner/clusters/clusterinfo"
	"testing"
)

func TestTransformAzureToHub(t *testing.T) {
	ci := &clusterinfo.ClusterInfo{
		Name:        "c",
		Cloud:       clusterinfo.AZURE,
		Location:    "westus2",
		Scope:       "joshua-playground",
		K8sVersion:  "1.14.0",
		GeneratedBy: clusterinfo.MOCK}
	tr := AKSTransformer{}
	std, err := tr.CloudToHub(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != "us-west1" {
		t.Error(std.Location)
	}
	if std.Cloud != clusterinfo.HUB {
		t.Errorf("not the standard cloud %s", std.Cloud)
	}
}
func TestTransformAzureToHubBadLoc(t *testing.T) {
	ci := &clusterinfo.ClusterInfo{Name: "c",
		Cloud: clusterinfo.AZURE, Location: "westus1",
		Scope:       "joshua-playground",
		K8sVersion:  "1.15.0",
		GeneratedBy: clusterinfo.MOCK}
	tr := AKSTransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}

func TestTransformHubToAzure(t *testing.T) {
	ci := &clusterinfo.ClusterInfo{
		Name:        "c",
		Cloud:       clusterinfo.HUB,
		Location:    "us-central1",
		Scope:       "",
		K8sVersion:  "1.14.0",
		GeneratedBy: clusterinfo.MOCK,
	}
	tr := AKSTransformer{}
	std, err := tr.HubToCloud(ci, "")
	if err != nil {
		t.Error(err)
	}
	if std.Location != "centralus" {
		t.Error(std.Location)
	}
	if std.Cloud != clusterinfo.AZURE {
		t.Errorf("Not the expected cloud %s", std.Cloud)
	}
}
