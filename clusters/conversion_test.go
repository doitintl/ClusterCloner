package clusters

import (
	"testing"
)

func TestTransformGcpToHubAndBack(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: GCP, Location: "us-east1", Scope: "joshua-playground"}
	std, err := tranformGCPToHub(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != ci.Location {
		t.Error(std.Location)
	}
	if std.Cloud != HUB {
		t.Errorf("Not the standard cloud %s", std.Cloud)
	}

	gcp, err := transformHubToGCP(std)
	if err != nil {
		t.Error(err)
	}
	if gcp.Scope != "" || gcp.Name != ci.Name || gcp.NodeCount != ci.NodeCount || gcp.Location != ci.Location || gcp.Cloud != ci.Cloud {
		t.Error(gcp)
	}
}
func TestTransformGcpToHubBadLoc(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: GCP, Location: "westus2", Scope: "joshua-playground"}
	_, err := tranformGCPToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}

func TestTransformAzureToHub(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: AZURE, Location: "westus2", Scope: "joshua-playground"}
	std, err := transformAzureToHub(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != "us-west1" {
		t.Error(std.Location)
	}
	if std.Cloud != HUB {
		t.Errorf("Not the standard cloud %s", std.Cloud)
	}
}
func TestTransformAzureToHubBadLoc(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: AZURE, Location: "westus1", Scope: "joshua-playground"}
	_, err := transformAzureToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}

func TestTransformHubToAzure(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: HUB, Location: "us-central1", Scope: ""}
	std, err := transformHubToAzure(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != "centralus" {
		t.Error(std.Location)
	}
	if std.Cloud != AZURE {
		t.Errorf("Not the expected cloud %s", std.Cloud)
	}
}
func TestTransformAzureToGCP(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: AZURE, Location: "westus2", Scope: "joshua-playground"}
	gcp, err := Transform(ci, GCP)
	if err != nil {
		t.Error(err)
	}
	if gcp.Location != "us-west1" {
		t.Error(gcp.Location)
	}
	if gcp.Cloud != GCP {
		t.Errorf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != "" || gcp.Name != ci.Name || gcp.NodeCount != ci.NodeCount || gcp.Location != "us-west1" {
		t.Error(gcp)
	}
}
