package access

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func init() {
	supportedVersions = []string{"1.14.9", "1.14.8", "1.14.11", "1.15.1", "1.15.8"}
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
func TestMachineTypes(t *testing.T) {
	types := MachineTypes
	log.Println(types)
	machineTypeCount := len(types)
	assert.Greater(t, machineTypeCount, 300)
	assert.Less(t, machineTypeCount, 330)
	mt := MachineTypeByName("Standard_A2_v2")
	assert.Equal(t, "Standard_A2_v2",mt.Name)
		assert.Equal(t, 1792,mt.RAMMB)
	assert.Equal(t, 1, mt.CPU )
}
