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
		MachineType: access.ParseMachineType(machineType),
		NodeCount:   1,
		K8sVersion:  "1.14.0",
		DiskSizeGB:  10}
	npi2 := clusterinfo.NodePoolInfo{
		Name:        "NP2",
		MachineType: access.ParseMachineType("Standard_A1"),
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
	m1 := clusterinfo.MachineType{Name: "m1-ultramem-80", CPU: 80, RAMGB: 1922}
	n1 := clusterinfo.MachineType{Name: "n1-ultramem-80", CPU: 80, RAMGB: 1922}
	if mtGcp != m1 && mtGcp != n1 {
		t.Error(mtGcp)
	}

}
