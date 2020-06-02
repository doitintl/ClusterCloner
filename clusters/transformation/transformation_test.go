package transformation

import (
	"clustercloner/clusters"
	accessaks "clustercloner/clusters/clouds/aks/access"
	accesseks "clustercloner/clusters/clouds/eks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
	"clustercloner/clusters/machinetypes"
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransformAzureToGCP(t *testing.T) {
	scope, azure := getSampleInputAKSCluster(t)
	gcp, err := transformCloudToCloud(azure, clusters.GCP, scope, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(gcp.Location, "us-west1") {
		t.Fatal(gcp.Location)
	}
	if gcp.Cloud != clusters.GCP {
		t.Fatalf("Not the right cloud %s", gcp.Cloud)
	}
	if gcp.Scope != scope ||
		gcp.Name != azure.Name ||
		!strings.HasPrefix(gcp.Location, "us-west1") ||
		len(gcp.NodePools) != len(azure.NodePools) {

		outputStr := util.ToJSON(gcp)
		inputStr := util.ToJSON(azure)
		t.Fatal(outputStr + "!=" + inputStr)
	}

	for i := range gcp.NodePools {
		azureNP := azure.NodePools[i]
		//Machine types and K8s versions will not match, so comparing NodePools with zeroed  Machine Types and K8s version
		azureNP.MachineType = machinetypes.MachineType{}
		azureNP.K8sVersion = ""
		gcpNP := gcp.NodePools[i]
		gcpNP.MachineType = machinetypes.MachineType{}
		gcpNP.K8sVersion = ""
		assert.Equal(t, gcpNP, azureNP)
	}

}

func getSampleInputAKSCluster(t *testing.T) (scope string, aksCluster *clusters.ClusterInfo) {
	scope = "sample-scope"
	machineType1ByName := accessaks.MachineTypeByName("Standard_D32s_v3")
	assert.NotEqual(t, machineType1ByName.Name, "")
	npi := clusters.NodePoolInfo{
		Name:        "NP",
		MachineType: machineType1ByName,
		NodeCount:   1,
		K8sVersion:  "1.14.0",
		DiskSizeGB:  10,
		Preemptible: true,
	}
	machineType2ByName := accessaks.MachineTypeByName("Standard_A2_v2")
	assert.NotEqual(t, machineType2ByName.Name, "")
	npi2 := clusters.NodePoolInfo{
		Name:        "NP2",
		MachineType: machineType2ByName,
		NodeCount:   2,
		K8sVersion:  "1.15.0",
		DiskSizeGB:  20,
		Preemptible: true,
	}

	npis := []clusters.NodePoolInfo{npi, npi2}
	nodePoolInfos := npis[:]

	aksCluster = &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.Azure,
		Location:    "westus2",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		NodePools:   nodePoolInfos,
		Labels:      map[string]string{"a": "aa", "b": "bb"},
		GeneratedBy: clusters.Mock,
	}
	return scope, aksCluster
}

func TestTransformGCPToAzure(t *testing.T) {
	gcpIn := sampleInputGCPCluster(t)
	azOut, err := transformCloudToCloud(gcpIn, clusters.Azure, gcpIn.Scope, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(azOut.Location, "centralus") {
		t.Fatal(azOut.Location)
	}
	if azOut.Cloud != clusters.Azure {
		t.Fatalf("Not the right cloud %s", azOut.Cloud)
	}
	if azOut.Scope != gcpIn.Scope ||
		azOut.Name != gcpIn.Name ||
		azOut.Location != "centralus" ||
		//			azOut.K8sVersion != gcpIn.K8sVersion ||
		len(azOut.NodePools) != len(gcpIn.NodePools) {
		outputStr := util.ToJSON(azOut)
		inputStr := util.ToJSON(gcpIn)
		t.Fatal(outputStr + "!=" + inputStr)
	}

	for i := range azOut.NodePools {
		//Zeroing out fields that are not expected to match
		npIn := gcpIn.NodePools[i]
		npIn.K8sVersion = ""
		npIn.MachineType = machinetypes.MachineType{}
		npOut := azOut.NodePools[i]
		npOut.MachineType = machinetypes.MachineType{}
		if !strings.HasPrefix(npOut.K8sVersion, "1.15") && !strings.HasPrefix(npOut.K8sVersion, "1.14") {
			t.Fatal(npOut.K8sVersion, "AKS may have upgraded versions")
		}
		npOut.K8sVersion = ""

		assert.Equal(t, npOut, npIn)
	}

	mtOut := azOut.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	expectedOutputMachineTypeNames := []string{
		"Standard_F16s",
		"Standard_F16",
		"Standard_F16s_v2",
	}
	found := false
	for _, mTypeName := range expectedOutputMachineTypeNames {
		expectedMachType := accessaks.MachineTypeByName(mTypeName)
		if expectedMachType == mtOut {
			found = true
		}
	}
	if !found {
		t.Fatal(mtOut.Name + " was not an expected machine type for " + gcpIn.NodePools[0].MachineType.Name + "; expected: " +
			strings.Join(expectedOutputMachineTypeNames, ","))
	}

}

func TestTransformGCPToAWS(t *testing.T) {
	gcpIn := sampleInputGCPCluster(t)
	awsOut, err := transformCloudToCloud(gcpIn, clusters.AWS, gcpIn.Scope, false)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "us-east-2", awsOut.Location)
	assert.Equal(t, clusters.AWS, awsOut.Cloud)

	if awsOut.Scope != gcpIn.Scope || awsOut.Name != gcpIn.Name || awsOut.Location != "us-east-2" || awsOut.K8sVersion != "1.14" || len(awsOut.NodePools) != len(gcpIn.NodePools) {
		t.Fatal(util.ToJSON(awsOut) + "!=" + util.ToJSON(gcpIn))
	}

	for i := range awsOut.NodePools {
		npIn := gcpIn.NodePools[i]
		npIn.K8sVersion = "" //Zeroing out fields that are not expected to match
		npIn.MachineType = machinetypes.MachineType{}
		npOut := awsOut.NodePools[i]
		npOut.MachineType = machinetypes.MachineType{}
		if npOut.K8sVersion != "1.15" && npOut.K8sVersion != "1.14" { //The test has different K8s Versions on different NodePools
			t.Fatal(npOut.K8sVersion)
		}
		npOut.K8sVersion = "" //Zeroing out fields that are not expected to match

		assert.Equal(t, npOut, npIn)
	}

	mtOut := awsOut.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	expectedOutputMachineTypeNames := []string{
		"c3.4xlarge",
		"c4.4xlarge",
	}
	if !util.ContainsStr(expectedOutputMachineTypeNames, mtOut.Name) {
		t.Fatal(mtOut.Name + " was not an expected machine type for " + gcpIn.NodePools[0].MachineType.Name)
	}

}

func TestTransformAWSToAzure(t *testing.T) {
	awsIn := sampleInputAWSCluster(t)
	azOut, err := transformCloudToCloud(awsIn, clusters.Azure, awsIn.Scope, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(azOut.Location, "centralus") {
		t.Fatal(azOut.Location)
	}
	if azOut.Cloud != clusters.Azure {
		t.Fatalf("Not the right cloud %s", azOut.Cloud)
	}
	if azOut.Scope != awsIn.Scope ||
		azOut.Name != awsIn.Name ||
		azOut.Location != "centralus" ||
		len(azOut.NodePools) != len(awsIn.NodePools) {
		outputStr := util.ToJSON(azOut)
		inputStr := util.ToJSON(awsIn)
		t.Fatal(outputStr + "!=" + inputStr)
	}

	for i := range azOut.NodePools {
		//Zeroing out fields that are not expected to match
		npIn := awsIn.NodePools[i]
		npIn.K8sVersion = ""
		npIn.MachineType = machinetypes.MachineType{}
		npOut := azOut.NodePools[i]
		npOut.MachineType = machinetypes.MachineType{}
		if !strings.HasPrefix(npOut.K8sVersion, "1.15") && !strings.HasPrefix(npOut.K8sVersion, "1.14") {
			t.Fatal(npOut.K8sVersion, "AKS may have upgraded versions")
		}
		npOut.K8sVersion = ""

		assert.Equal(t, npOut, npIn)
	}

	mtOut := azOut.NodePools[0].MachineType

	// Can vary because map is not deterministically ordered
	expectedOutputMachineTypeNames := []string{
		"Standard_F16s",
		"Standard_F16",
		"Standard_F16s_v2",
		"Standard_F8s",
		"Standard_F8s_v2",
		"Standard_A8_v2",
	}
	if !util.ContainsStr(expectedOutputMachineTypeNames, mtOut.Name) {
		t.Fatal(mtOut.Name + " was not an expected machine type for " + awsIn.NodePools[0].MachineType.Name + "; expected: " +
			strings.Join(expectedOutputMachineTypeNames, ","))
	}

}

func sampleInputAWSCluster(t *testing.T) *clusters.ClusterInfo {

	scope := "sample-project"
	inputMachTypeFirstNode := "c3.2xlarge"
	machTypeByName := accesseks.MachineTypeByName(inputMachTypeFirstNode)
	if machTypeByName.Name == "" {
		t.Fatal("cannot find machine type " + inputMachTypeFirstNode)
	}
	npi1 := clusters.NodePoolInfo{
		Name:        "NP",
		MachineType: machTypeByName,
		NodeCount:   1,
		K8sVersion:  "1.14",
		DiskSizeGB:  10,
		Preemptible: true,
	}
	machTypeByName2 := accesseks.MachineTypeByName("c3.4xlarge")
	if machTypeByName2.Name == "" {
		t.Fatal("cannot find machine type")
	}
	npi2 := clusters.NodePoolInfo{
		Name:        "NP2",
		MachineType: machTypeByName2,
		NodeCount:   2,
		K8sVersion:  "1.15",
		DiskSizeGB:  20,
		Preemptible: true,
	}

	npis := []clusters.NodePoolInfo{npi1, npi2}
	nodePools := npis[:]

	awsIn := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.AWS,
		Location:    "us-east-2",
		Scope:       scope,
		K8sVersion:  "1.14",
		Labels:      map[string]string{"a": "aa", "b": "bb"},
		NodePools:   nodePools,
		GeneratedBy: clusters.Mock,
	}
	return awsIn
}

func sampleInputGCPCluster(t *testing.T) *clusters.ClusterInfo {
	scope := "sample-project"
	inputMachTypeFirstNode := "e2-highcpu-16"
	machTypeByName1 := accessgke.MachineTypeByName(inputMachTypeFirstNode)
	if machTypeByName1.Name == "" {
		t.Fatal("cannot find machine type")
	}
	npi1 := clusters.NodePoolInfo{
		Name:        "NP",
		MachineType: machTypeByName1,
		NodeCount:   1,
		K8sVersion:  "1.14.3",
		DiskSizeGB:  10,
		Preemptible: true,
	}
	machTypeByName2 := accessgke.MachineTypeByName("c2-standard-60")
	if machTypeByName2.Name == "" {
		t.Fatal("cannot find machine type")
	}
	npi2 := clusters.NodePoolInfo{
		Name:        "NP2",
		MachineType: machTypeByName2,
		NodeCount:   2,
		K8sVersion:  "1.15.2",
		DiskSizeGB:  20,
		Preemptible: true,
	}

	npis := []clusters.NodePoolInfo{npi1, npi2}
	nodePools := npis[:]

	sampleGCPClsuter := &clusters.ClusterInfo{
		Name:        "c",
		Cloud:       clusters.GCP,
		Location:    "us-central1-c",
		Scope:       scope,
		K8sVersion:  "1.14.0",
		Labels:      map[string]string{"a": "aa", "b": "bb"},
		NodePools:   nodePools,
		GeneratedBy: clusters.Mock,
	}
	return sampleGCPClsuter
}
