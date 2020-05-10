package nodes

import (
	accessaks "clustercloner/clusters/clouds/aks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
	"testing"
)

func TestFindMatchingMachineTypeGkeToAks(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt := accessgke.MachineTypeByName(gkeMachine)
	matching := FindMatchingMachineType(mt, accessaks.MachineTypesNoPromo)
	if matching.CPU != 8 || matching.RAMGB != 14 {
		t.Errorf("No match: %v", matching)
	}
}

//not realistic use of
func TestFindMatchingMachineTypeGkeToGke(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt := accessgke.MachineTypeByName(gkeMachine)
	matching := FindMatchingMachineType(mt, accessgke.MachineTypes)
	if matching.CPU != 8 || matching.RAMGB != 8 {
		t.Errorf("No match: %v", matching)
	}
}
