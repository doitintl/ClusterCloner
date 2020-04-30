package clusters

import (
	"clusterCloner/clusters/cluster_conversion"
	"clusterCloner/clusters/cluster_info"
	"testing"
)

func TestTransformAzureToGCP(t *testing.T) {
	ci := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.AZURE, Location: "westus2", Scope: "joshua-playground"}
	gcp, err := cluster_conversion.Transform(ci, cluster_info.GCP)
	if err != nil {
		t.Error(err)
	}
	if gcp.Location != "us-west1" {
		t.Error(gcp.Location)
	}
	if gcp.Cloud != cluster_info.GCP {
		t.Errorf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != "" || gcp.Name != ci.Name || gcp.NodeCount != ci.NodeCount || gcp.Location != "us-west1" {
		t.Error(gcp)
	}
}
