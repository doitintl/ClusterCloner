package nodes

import (
	accessaks "clustercloner/clusters/clouds/aks/access"
	accessgke "clustercloner/clusters/clouds/gke/access"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindMatchingMachineTypeGkeToAks(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt, err := accessgke.GetMachineTypes().Get(gkeMachine)
	assert.Nil(t, err)
	matching := FindMatchingMachineType(mt, accessaks.GetMachineTypes())
	if matching.CPU != 8 || matching.RAMMB != 16384 {
		t.Fatalf("No match: %v", matching)
	}
}

//not realistic use of
func TestFindMatchingMachineTypeGkeToGke(t *testing.T) {
	gkeMachine := "n2d-highcpu-8"
	mt, err := accessgke.GetMachineTypes().Get(gkeMachine)
	assert.Nil(t, err)
	matching := FindMatchingMachineType(mt, accessgke.GetMachineTypes())
	if matching.CPU != 8 || matching.RAMMB != 8000 {
		t.Fatalf("No match: %v", matching)
	}
}
