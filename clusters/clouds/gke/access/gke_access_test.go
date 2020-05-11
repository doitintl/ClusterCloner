package access

import (
	"clustercloner/clusters"
	"testing"
)

func TestParseMachineType(t *testing.T) {
	machineType := "e2-highcpu-8"
	mt := MachineTypeByName(machineType)
	if mt.Name != machineType {
		t.Error(mt.Name)
	}
	if mt.CPU != 8 {
		t.Error(mt.CPU)

	}
	if mt.RAMGB != 8 {
		t.Error(mt.RAMGB)
	}
}
func TestParseMachineType2(t *testing.T) {
	name := "n1-ultramem-40"
	mt := MachineTypeByName(name)
	if mt.Name != name {
		t.Error(mt.Name)
	}
	if mt.CPU != 40 {
		t.Error(mt.CPU)
	}
	if mt.RAMGB != 961 {
		t.Error(mt.RAMGB)
	}
}
func TestParseMissingMachineType2(t *testing.T) {
	name := "xx-xx-40"
	mt := MachineTypeByName(name)
	zero := clusters.MachineType{}
	if mt != zero {
		t.Error(mt.Name)
	}
}
