package transformation

import (
	"clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransformAzureToGCP(t *testing.T) {
	scope := "joshua-playground"
	machineType := "Standard_M64ms"
	npi := clusterinfo.NodePoolInfo{
		Name:        "NP",
		MachineType: access.MachineTypeByName(machineType),
		NodeCount:   1,
		K8sVersion:  "1.14.0",
		DiskSizeGB:  10}
	npi2 := clusterinfo.NodePoolInfo{
		Name:        "NP2",
		MachineType: access.MachineTypeByName("Standard_A1"),
		NodeCount:   2,
		K8sVersion:  "1.15.0",
		DiskSizeGB:  20}

	npis := []clusterinfo.NodePoolInfo{npi, npi2}
	nodePools := npis[:]

	azure := &clusterinfo.ClusterInfo{
		Name:        "c",
		Cloud:       clusterinfo.AZURE,
		Location:    "westus2",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		NodePools:   nodePools,
		GeneratedBy: clusterinfo.MOCK,
	}
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
	if gcp.Scope != scope ||
		gcp.Name != azure.Name ||
		!strings.HasPrefix(gcp.Location, "us-west1") ||
		gcp.K8sVersion != azure.K8sVersion ||
		len(gcp.NodePools) != len(azure.NodePools) {
		outputStr := util.MarshallToJSONString(gcp)
		inputStr := util.MarshallToJSONString(azure)
		t.Error(outputStr + "!=" + inputStr)
	}

	for i := range gcp.NodePools {
		azureNP := azure.NodePools[i]
		//Machine types will not match, so comparing NodePools with zero Machine Types
		azureNP.MachineType = clusterinfo.MachineType{}
		gcpNP := gcp.NodePools[i]
		gcpNP.MachineType = clusterinfo.MachineType{}
		assert.Equal(t, gcpNP, azureNP)
	}

	mtGcp := gcp.NodePools[0].MachineType

	// Can vary because map is not determinstically ordered
	machineA := clusterinfo.MachineType{Name: "m1-ultramem-80", CPU: 80, RAMGB: 1922}
	machineB := clusterinfo.MachineType{Name: "n1-ultramem-80", CPU: 80, RAMGB: 1922}
	if mtGcp != machineA && mtGcp != machineB {
		t.Error(mtGcp)
	}
}

func TestTransformGCPToAzure(t *testing.T) {
	scope := "joshua-playground"
	machineType := "e2-highcpu-16"
	npi := clusterinfo.NodePoolInfo{
		Name:        "NP",
		MachineType: access.MachineTypeByName(machineType),
		NodeCount:   1,
		K8sVersion:  "1.14.3",
		DiskSizeGB:  10}
	npi2 := clusterinfo.NodePoolInfo{
		Name:        "NP2",
		MachineType: access.MachineTypeByName("Standard_A1"),
		NodeCount:   2,
		K8sVersion:  "1.15.2",
		DiskSizeGB:  20}

	npis := []clusterinfo.NodePoolInfo{npi, npi2}
	nodePools := npis[:]

	gcpIn := &clusterinfo.ClusterInfo{
		Name:        "c",
		Cloud:       clusterinfo.GCP,
		Location:    "us-central1-c",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		NodePools:   nodePools,
		GeneratedBy: clusterinfo.MOCK,
	}
	azOut, err := transformCloudToCloud(gcpIn, clusterinfo.AZURE, scope)
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(azOut.Location, "centralus") {
		t.Error(azOut.Location)
	}
	if azOut.Cloud != clusterinfo.AZURE {
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
		npIn.MachineType = clusterinfo.MachineType{}
		npOut := azOut.NodePools[i]
		npOut.MachineType = clusterinfo.MachineType{}
		if npOut.K8sVersion != "1.15.7" && npOut.K8sVersion != "1.14.7" {
			t.Error(npOut.K8sVersion)
		}
		npOut.K8sVersion = ""

		assert.Equal(t, npOut, npIn)
	}

	mtOut := azOut.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	mTypeNames := []string{"Standard_B1s",
		"Basic_A1",
		"Basic_A0",
		"Standard_A1",
		"Standard_A0",
		"Standard_B1ls",
	}
	found := false
	for _, mTypeName := range mTypeNames {
		mType := access.MachineTypeByName(mTypeName)
		if mType == mtOut {
			found = true
		}
	}
	if !found {
		t.Error("Cannot find " + mtOut.Name)
	}

}
