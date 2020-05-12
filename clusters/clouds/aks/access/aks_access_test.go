package access

import (
	"clustercloner/clusters"
	"log"
	"testing"
)

func init() {
	supportedVersions = []string{"1.14.9", "1.14.8", "1.14.11", "1.15.1", "1.15.8"}
}

func TestDescribeCluster(t *testing.T) {
	a := AKSClusterAccess{}
	ci := clusters.ClusterInfo{
		Scope:    "joshua-playground",
		Location: "westus2",
		Name:     "cluster-2",
	}
	ciRead, _ := a.DescribeCluster(&ci)
	log.Println(ciRead)
	log.Println(ciRead)
}
func TestParseMachineType(t *testing.T) {
	machineType := "Standard_D1_v2"

	mt := MachineTypeByName(machineType)
	if mt.Name != machineType {
		t.Error(mt.Name)
	}
	if mt.CPU != 1 {
		t.Error(mt.CPU)

	}
	if mt.RAMMB != 3584 {
		t.Error(mt.RAMMB)
	}
}
