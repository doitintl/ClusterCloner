package access

import (
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
	if mt.RAMGB != 3 {
		t.Error(mt.RAMGB)
	}
}
