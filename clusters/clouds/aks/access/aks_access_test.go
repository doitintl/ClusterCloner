package access

import (
	"testing"
)

func TestParseMachineType(t *testing.T) {
	machineType := "Standard_D1_v2"

	mt := ParseMachineType(machineType)
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
