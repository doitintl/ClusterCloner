package transform

import (
	"clusterCloner/clusters/cluster_info"
	"testing"
)

func TestTransformAzureToHub(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.AZURE, Location: "westus2", Scope: "joshua-playground", GeneratedBy: cluster_info.MOCK}
	tr := AksTransformer{}
	std, err := tr.CloudToHub(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != "us-west1" {
		t.Error(std.Location)
	}
	if std.Cloud != cluster_info.HUB {
		t.Errorf("Not the standard cloud %s", std.Cloud)
	}
}
func TestTransformAzureToHubBadLoc(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.AZURE, Location: "westus1", Scope: "joshua-playground", GeneratedBy: cluster_info.MOCK}
	tr := AksTransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}

func TestTransformHubToAzure(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.HUB, Location: "us-central1", Scope: "", GeneratedBy: cluster_info.MOCK}
	tr := AksTransformer{}
	std, err := tr.HubToCloud(ci, "")
	if err != nil {
		t.Error(err)
	}
	if std.Location != "centralus" {
		t.Error(std.Location)
	}
	if std.Cloud != cluster_info.AZURE {
		t.Errorf("Not the expected cloud %s", std.Cloud)
	}
}
