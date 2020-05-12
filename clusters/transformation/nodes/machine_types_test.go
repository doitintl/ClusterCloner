package nodes

import (
	accessaks "clustercloner/clusters/clouds/aks/access"
	"clustercloner/clusters/clouds/gke/access"
	"testing"
)

func TestFindMatchingMachineTypeGkeToAks(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt := access.MachineTypeByName(gkeMachine)
	matching := FindMatchingMachineType(mt, accessaks.MachineTypes)
	if matching.CPU != 8 || matching.RAMMB != 16384 {
		t.Errorf("No match: %v", matching)
	}
}

//not realistic use of
func TestFindMatchingMachineTypeGkeToGke(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt := access.MachineTypeByName(gkeMachine)
	matching := FindMatchingMachineType(mt, access.MachineTypes)
	if matching.CPU != 8 || matching.RAMMB != 8000 {
		t.Errorf("No match: %v", matching)
	}
}
