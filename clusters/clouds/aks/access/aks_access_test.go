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

func TestSupportedK8sVersion(t *testing.T) {

	matchingSupported, err := FindBestMatchingSupportedK8sVersion("1.14.1")
	if err != nil {
		t.Error(err)
	}
	if matchingSupported != "1.14.8" {
		t.Error(matchingSupported)
	}
}

func TestSupportedK8sVersionError(t *testing.T) {

	_, err := FindBestMatchingSupportedK8sVersion("1.214.10")
	if err == nil {
		t.Error(err)
	}

}

func TestSupportedK8sVersion3(t *testing.T) {

	supported, err := FindBestMatchingSupportedK8sVersion("1.14.10")
	if err != nil {
		t.Error(err)
	}
	if supported != "1.14.11" {
		t.Error(supported)
	}
}
