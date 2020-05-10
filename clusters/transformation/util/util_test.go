package util

import (
	"clustercloner/clusters/transformation/nodes/util"
	"testing"
)

func TestMajorMinorPatchVersion(t *testing.T) {
	s, e := util.MajorMinorPatchVersion("1.14.10-gke.27")
	if e != nil || s != "1.14.10" {
		t.Errorf("No match: %s", s)
	}
}
func TestMajorMinorPatchVersionOnly(t *testing.T) {
	s, e := util.MajorMinorPatchVersion("1.14.10")
	if e != nil || s != "1.14.10" {
		t.Errorf("No match: %s", s)
	}
}

func TestMajorMinorVersionIn(t *testing.T) {
	s, e := util.MajorMinorPatchVersion("1.14")
	if e != nil || s != "1.14.0" {
		t.Errorf("No match: %s", s)
	}
}

func TestMajorMinorPatchVersionFail(t *testing.T) {
	s, e := util.MajorMinorPatchVersion("11410-gke.27")
	if e == nil {
		t.Errorf("Expect match: %s", s)
	}
}
func TestMajorMinorVersionOut(t *testing.T) {
	s, e := util.MajorMinorVersion("1.14.10-gke.27")
	if e != nil || s != "1.14" {
		t.Errorf("No match: %s", s)
	}
}
func TestMajorMinorVersionInOut(t *testing.T) {
	s, e := util.MajorMinorVersion("1.14")
	if e != nil || s != "1.14" {
		t.Errorf("No match: %s", s)
	}
}
func TestMajorMinorVersionOut2(t *testing.T) {
	s, e := util.MajorMinorVersion("1.14.10-gke.27")
	if e != nil || s != "1.14" {
		t.Errorf("No match: %s", s)
	}
}

func TestPatchVersion(t *testing.T) {
	s, e := util.PatchVersion("1.14.10-gke.27")
	if e != nil || s != 10 {
		t.Errorf("No match: %d", s)
	}
}

func TestPatchVersion2(t *testing.T) {
	s, e := util.PatchVersion("1.14.10")
	if e != nil || s != 10 {
		t.Errorf("No match: %d", s)
	}
}

func TestPatchVersionErr(t *testing.T) {
	s, e := util.PatchVersion("1.14")
	if e == nil {
		t.Errorf("Did not expect match: %d", s)
	}
}
