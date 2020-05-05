package util

import "testing"

func TestMajorMinorPatchVersion(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14.10-gke.27")
	if e != nil || s != "1.14.10" {
		t.Errorf("No match: %s", s)
	}
}
func TestMajorMinorPatchVersionOnly(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14.10")
	if e != nil || s != "1.14.10" {
		t.Errorf("No match: %s", s)
	}
}

func TestMajorMinorVersion(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14")
	if e != nil || s != "1.14.0" {
		t.Errorf("No match: %s", s)
	}
}

func TestMajorMinorPatchVersionFail(t *testing.T) {
	s, e := MajorMinorPatchVersion("11410-gke.27")
	if e == nil {
		t.Errorf("Expect match: %s", s)
	}
}
