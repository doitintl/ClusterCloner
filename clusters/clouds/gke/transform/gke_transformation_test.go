package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/gke/access"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransformGcpToHubAndBack(t *testing.T) {
	scope := "sample-project"
	mt, err := access.MachineTypes.Get("e2-highcpu-8")
	assert.Nil(t, err)
	var npi1 = clusters.NodePoolInfo{
		Name:        "NPName",
		NodeCount:   2,
		K8sVersion:  "1.15.2-gke27",
		MachineType: mt,
		DiskSizeGB:  32,
		Preemptible: true,
	}
	npi2 := npi1 //copy
	npi2.Name = "yyy"
	nodePools := [2]clusters.NodePoolInfo{npi1, npi2}
	input := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.GCP,
		Location:    "us-east1-a",
		K8sVersion:  "1.14.1-gke27",
		Scope:       scope,
		GeneratedBy: clusters.Mock,
		Labels:      map[string]string{"x": "y"},
		NodePools:   nodePools[:],
	}
	tr := GKETransformer{}
	hub, err := tr.CloudToHub(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(input.Location, hub.Location) {
		t.Fatal(hub.Location)
	}
	if hub.Cloud != clusters.Hub {
		t.Fatalf("Not the Hub: %s", hub.Cloud)
	}

	output, err := tr.HubToCloud(hub, scope)
	if err != nil {
		t.Fatal(err)
	}

	if output.Scope != scope || output.Name != input.Name ||
		output.Cloud != input.Cloud {
		outputStr := util.ToJSON(output)
		inputStr := util.ToJSON(input)
		t.Fatal(outputStr + "!=" + inputStr)
	}
	if output.NodePools[0].DiskSizeGB != input.NodePools[0].DiskSizeGB {
		t.Fatal(output.NodePools[0].DiskSizeGB)
	}
}
func TestTransformGcpToHubBadLoc(t *testing.T) {
	ci := &clusters.ClusterInfo{Name: "c",
		Cloud:       clusters.GCP,
		Location:    "westus2",
		K8sVersion:  "1.14.1-gke-27",
		Scope:       "sample-project",
		Labels:      map[string]string{"x": "y"},
		GeneratedBy: clusters.Mock,
	}
	tr := GKETransformer{}
	_, err := tr.CloudToHub(ci)
	if err == nil {
		t.Fatal("expect error")
	}
}
func TestHyphens(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("us-central1-c")
	if hyCount != 2 {
		t.Fatal(hyCount)
	}
	if secondHyIdx != 11 {
		t.Fatal(secondHyIdx)
	}
}
func TestHyphensNone(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("uscentral1c")
	if hyCount != 0 {
		t.Fatal(hyCount)
	}
	if secondHyIdx != -1 {
		t.Fatal(secondHyIdx)
	}
}
func TestHyphenOne(t *testing.T) {
	hyCount, secondHyIdx := hyphensForGCPLocation("us-central1")
	if hyCount != 1 {
		t.Fatal(hyCount)
	}
	if secondHyIdx != -1 {
		t.Fatal(secondHyIdx)
	}
}
