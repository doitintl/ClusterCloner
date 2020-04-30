package transform

import (
	"clusterCloner/clusters/cluster_info"
	"testing"
)

func TestTransformGcpToHubAndBack(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.GCP, Location: "us-east1", Scope: "joshua-playground"}
	std, err := TranformGCPToHub(ci)
	if err != nil {
		t.Error(err)
	}
	if std.Location != ci.Location {
		t.Error(std.Location)
	}
	if std.Cloud != cluster_info.HUB {
		t.Errorf("Not the standard cloud %s", std.Cloud)
	}

	gcp, err := TransformHubToGCP(std)
	if err != nil {
		t.Error(err)
	}
	if gcp.Scope != "" || gcp.Name != ci.Name || gcp.NodeCount != ci.NodeCount || gcp.Location != ci.Location || gcp.Cloud != ci.Cloud {
		t.Error(gcp)
	}
}
func TestTransformGcpToHubBadLoc(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.GCP, Location: "westus2", Scope: "joshua-playground"}
	_, err := TranformGCPToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}
