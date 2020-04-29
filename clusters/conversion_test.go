package clusters

import "testing"

func TestConvertGcpToHub(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: GCP, Location: "us-east1", Scope: "joshua-playground"}
	std, err := ConvertGCPToStandard(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != ci.Location {
		t.Error(std.Location)
	}
	if std.Cloud != HUB {
		t.Errorf("Not the standard cloud %s", std.Cloud)
	}
}
func TestConvertGcpToHubBadLoc(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: GCP, Location: "westus2", Scope: "joshua-playground"}
	_, err := ConvertGCPToStandard(ci)
	if err == nil {
		t.Error("expect error")
	}
}

func TestConvertAzureToHub(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: AZURE, Location: "westus2", Scope: "joshua-playground"}
	std, err := ConvertAzureToStandard(ci)
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
func TestConvertAzureToHubBadLoc(t *testing.T) {
	ci := ClusterInfo{Name: "c", NodeCount: 1, Cloud: AZURE, Location: "westus1", Scope: "joshua-playground"}
	_, err := ConvertAzureToStandard(ci)
	if err == nil {
		t.Error("expect error")
	}
}
