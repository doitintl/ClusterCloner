package transformation

import (
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/util"
	"log"
	"strings"
	"testing"
)

func TestTransformAzureToGCP(t *testing.T) {
	scope := "joshua-playground"
	azure := clusterinfo.ClusterInfo{
		Name:        "c",
		NodeCount:   1,
		Cloud:       clusterinfo.AZURE,
		Location:    "westus2",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		GeneratedBy: clusterinfo.MOCK}
	gcp, err := transformCloudToCloud(azure, clusterinfo.GCP, scope)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(gcp.Location, "us-west1") {
		t.Error(gcp.Location)
	}
	if gcp.Cloud != clusterinfo.GCP {
		t.Errorf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != scope || gcp.Name != azure.Name || gcp.NodeCount != azure.NodeCount || !strings.HasPrefix(gcp.Location, "us-west1") {
		outputStr := util.MarshallToJSONString(gcp)
		inputStr := util.MarshallToJSONString(azure)
		t.Error(outputStr + "!=" + inputStr)
	}
	log.Println(gcp.K8sVersion)
	log.Println(azure.K8sVersion)
}
