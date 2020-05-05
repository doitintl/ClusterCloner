package transform

import (
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"strings"
	"testing"
)

func TestTransformGcpToHubAndBack(t *testing.T) {
	scope := "joshua-playground"
	input := cluster_info.ClusterInfo{
		Name:        "c",
		NodeCount:   1,
		Cloud:       cluster_info.GCP,
		Location:    "us-east1-a",
		K8sVersion:  "1.14.1-gke27",
		Scope:       scope,
		GeneratedBy: cluster_info.MOCK}
	tr := GkeTransformer{}
	hub, err := tr.CloudToHub(input)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(input.Location, hub.Location) {
		t.Error(hub.Location)
	}
	if hub.Cloud != cluster_info.HUB {
		t.Errorf("Not the standard cloud %s", hub.Cloud)
	}

	output, err := tr.HubToCloud(hub, scope)
	if err != nil {
		t.Error(err)
	}

	if output.Scope != scope || output.Name != input.Name || output.NodeCount != input.NodeCount ||
		output.Cloud != input.Cloud {
		outputStr := util.MarshallToJsonString(output)
		inputStr := util.MarshallToJsonString(input)
		t.Error(outputStr + "!=" + inputStr)
	}
}
func TestTransformGcpToHubBadLoc(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.GCP, Location: "westus2",
		K8sVersion: "1.14.1-gke-27",
		Scope:      "joshua-playground", GeneratedBy: cluster_info.MOCK}
	tr := GkeTransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}
