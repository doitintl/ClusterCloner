package transform

import (
	"clustercloner/clusters/clouds/gke/access"
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/util"
	"strings"
	"testing"
)

func TestTransformGcpToHubAndBack(t *testing.T) {
	scope := "joshua-playground"
	mt := access.ParseMachineType("e2-highcpu-8")
	var npi1 = clusterinfo.NodePoolInfo{
		Name:        "xx",
		NodeCount:   2,
		K8sVersion:  "1/15.2-gke27",
		MachineType: mt,
		DiskSizeGB:  32,
	}
	npi2 := npi1 //copy
	npi2.Name = "yyy"
	nodePools := [2]clusterinfo.NodePoolInfo{npi1, npi2}
	input := clusterinfo.ClusterInfo{
		Name:                "c",
		DeprecatedNodeCount: 1,
		Cloud:               clusterinfo.GCP,
		Location:            "us-east1-a",
		K8sVersion:          "1.14.1-gke27",
		Scope:               scope,
		GeneratedBy:         clusterinfo.MOCK,
		NodePools:           nodePools[:],
	}
	tr := GKETransformer{}
	hub, err := tr.CloudToHub(input)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(input.Location, hub.Location) {
		t.Error(hub.Location)
	}
	if hub.Cloud != clusterinfo.HUB {
		t.Errorf("Not the standard cloud %s", hub.Cloud)
	}

	output, err := tr.HubToCloud(hub, scope)
	if err != nil {
		t.Error(err)
	}

	if output.Scope != scope || output.Name != input.Name || output.DeprecatedNodeCount != input.DeprecatedNodeCount ||
		output.Cloud != input.Cloud {
		outputStr := util.MarshallToJSONString(output)
		inputStr := util.MarshallToJSONString(input)
		t.Error(outputStr + "!=" + inputStr)
	}
	if output.NodePools[0].DiskSizeGB != input.NodePools[0].DiskSizeGB {
		t.Error(output.NodePools[0].DiskSizeGB)
	}
}
func TestTransformGcpToHubBadLoc(t *testing.T) {
	ci := clusterinfo.ClusterInfo{Name: "c", DeprecatedNodeCount: 1, Cloud: clusterinfo.GCP, Location: "westus2",
		K8sVersion: "1.14.1-gke-27",
		Scope:      "joshua-playground", GeneratedBy: clusterinfo.MOCK}
	tr := GKETransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Error("expect error")
	}
}
func TestHyphens(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("us-central1-c")
	if hyCount != 2 {
		t.Error(hyCount)
	}
	if secondHyIdx != 11 {
		t.Error(secondHyIdx)
	}
}
func TestHyphensNone(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("uscentral1c")
	if hyCount != 0 {
		t.Error(hyCount)
	}
	if secondHyIdx != -1 {
		t.Error(secondHyIdx)
	}
}
func TestHyphenOne(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("us-central1")
	if hyCount != 1 {
		t.Error(hyCount)
	}
	if secondHyIdx != -1 {
		t.Error(secondHyIdx)
	}
}
