package cluster_transformation

import (
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"strings"
	"testing"
)

func TestTransformAzureToGCP(t *testing.T) {
	scope := "joshua-playground"
	azure := cluster_info.ClusterInfo{Name: "c", NodeCount: 1, Cloud: cluster_info.AZURE, Location: "westus2", Scope: scope, GeneratedBy: cluster_info.MOCK}
	gcp, err := transform(azure, cluster_info.GCP, scope)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(gcp.Location, "us-west1") {
		t.Error(gcp.Location)
	}
	if gcp.Cloud != cluster_info.GCP {
		t.Errorf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != scope || gcp.Name != azure.Name || gcp.NodeCount != azure.NodeCount || !strings.HasPrefix(gcp.Location, "us-west1") {
		outputStr := util.MarshallToJsonString(gcp)
		inputStr := util.MarshallToJsonString(azure)
		t.Error(outputStr + "!=" + inputStr)
	}
}
