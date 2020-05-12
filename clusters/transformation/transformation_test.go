package transformation

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clouds/gke/access"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TODO integration test with actual creation of clones
func TestTransformAzureToGCP(t *testing.T) {
	scope := "sample-scope"
	machineType := "Standard_M64ms"
	npi := clusters.NodePoolInfo{
		Name:        "NP",
		MachineType: accessaks.MachineTypeByName(machineType),
		NodeCount:   1,
		K8sVersion:  "1.14.0",
		DiskSizeGB:  10}
	npi2 := clusters.NodePoolInfo{
		Name:        "NP2",
		MachineType: accessaks.MachineTypeByName("Standard_A1"),
		NodeCount:   2,
		K8sVersion:  "1.15.0",
		DiskSizeGB:  20}

	npis := []clusters.NodePoolInfo{npi, npi2}
	nodePools := npis[:]

	azure := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.Azure,
		Location:    "westus2",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		NodePools:   nodePools,
		GeneratedBy: clusters.Mock,
	}
	gcp, err := transformCloudToCloud(azure, clusters.GCP, scope, false)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(gcp.Location, "us-west1") {
		t.Error(gcp.Location)
	}
	if gcp.Cloud != clusters.GCP {
		t.Errorf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != scope ||
		gcp.Name != azure.Name ||
		!strings.HasPrefix(gcp.Location, "us-west1") ||
		len(gcp.NodePools) != len(azure.NodePools) {

		outputStr := util.MarshallToJSONString(gcp)
		inputStr := util.MarshallToJSONString(azure)
		t.Error(outputStr + "!=" + inputStr)
	}

	for i := range gcp.NodePools {
		azureNP := azure.NodePools[i]
		//Machine types and K8s versions will not match, so comparing NodePools with zeroed  Machine Types and K8s version
		azureNP.MachineType = clusters.MachineType{}
		azureNP.K8sVersion = ""
		gcpNP := gcp.NodePools[i]
		gcpNP.MachineType = clusters.MachineType{}
		gcpNP.K8sVersion = ""
		assert.Equal(t, gcpNP, azureNP)
	}

	mtGcp := gcp.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	machineA := access.MachineTypeByName("f1-micro")
	if mtGcp != machineA {
		t.Error(mtGcp)
	}
}

func TestTransformGCPToAzure(t *testing.T) {
	scope := "sample-project"
	machineType := "e2-highcpu-16"
	npi := clusters.NodePoolInfo{
		Name:        "NP",
		MachineType: accessaks.MachineTypeByName(machineType),
		NodeCount:   1,
		K8sVersion:  "1.14.3",
		DiskSizeGB:  10}
	npi2 := clusters.NodePoolInfo{
		Name:        "NP2",
		MachineType: accessaks.MachineTypeByName("Standard_A1"),
		NodeCount:   2,
		K8sVersion:  "1.15.2",
		DiskSizeGB:  20}

	npis := []clusters.NodePoolInfo{npi, npi2}
	nodePools := npis[:]

	gcpIn := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.GCP,
		Location:    "us-central1-c",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		NodePools:   nodePools,
		GeneratedBy: clusters.Mock,
	}
	azOut, err := transformCloudToCloud(gcpIn, clusters.Azure, scope, false)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(azOut.Location, "centralus") {
		t.Error(azOut.Location)
	}
	if azOut.Cloud != clusters.Azure {
		t.Errorf("Not the right cloud %s", azOut.Cloud)
	}
	if azOut.Scope != scope ||
		azOut.Name != gcpIn.Name ||
		azOut.Location != "centralus" ||
		//			azOut.K8sVersion != gcpIn.K8sVersion ||
		len(azOut.NodePools) != len(gcpIn.NodePools) {
		outputStr := util.MarshallToJSONString(azOut)
		inputStr := util.MarshallToJSONString(gcpIn)
		t.Error(outputStr + "!=" + inputStr)
	}

	for i := range azOut.NodePools {
		//Zeroing out fields that are not expected to match
		npIn := gcpIn.NodePools[i]
		npIn.K8sVersion = ""
		npIn.MachineType = clusters.MachineType{}
		npOut := azOut.NodePools[i]
		npOut.MachineType = clusters.MachineType{}
		if npOut.K8sVersion != "1.15.7" && npOut.K8sVersion != "1.14.7" {
			t.Error(npOut.K8sVersion)
		}
		npOut.K8sVersion = ""

		assert.Equal(t, npOut, npIn)
	}

	mtOut := azOut.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	mTypeNames := []string{

		"Standard_DS1",
		"Standard_DS1_v2",
		"Standard_D1_v2",
	}
	found := false
	for _, mTypeName := range mTypeNames {
		mType := accessaks.MachineTypeByName(mTypeName)
		if mType == mtOut {
			found = true
		}
	}
	if !found {
		t.Error("Cannot find " + mtOut.Name)
	}

}
