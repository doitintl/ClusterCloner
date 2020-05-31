package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMajorMinorPatchVersion(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14.10-gke.27")
	if e != nil || s != "1.14.10" {
		t.Fatalf("No match: %s", s)
	}
}
func TestMajorMinorPatchVersionOnly(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14.10")
	if e != nil || s != "1.14.10" {
		t.Fatalf("No match: %s", s)
	}
}

func TestMajorMinorVersionIn(t *testing.T) {
	s, e := MajorMinorPatchVersion("1.14")
	if e != nil || s != "1.14.0" {
		t.Fatalf("No match: %s", s)
	}
}

func TestMajorMinorPatchVersionFail(t *testing.T) {
	s, e := MajorMinorPatchVersion("11410-gke.27")
	if e == nil {
		t.Fatalf("Expect match: %s", s)
	}
}
func TestMajorMinorVersionOut(t *testing.T) {
	s, e := MajorMinorVersion("1.14.10-gke.27")
	if e != nil || s != "1.14" {
		t.Fatalf("No match: %s", s)
	}
}
func TestMajorMinorVersionInOut(t *testing.T) {
	s, e := MajorMinorVersion("1.14")
	if e != nil || s != "1.14" {
		t.Fatalf("No match: %s", s)
	}
}
func TestMajorMinorVersionOut2(t *testing.T) {
	s, e := MajorMinorVersion("1.14.10-gke.27")
	if e != nil || s != "1.14" {
		t.Fatalf("No match: %s", s)
	}
}

func TestPatchVersion(t *testing.T) {
	s, e := PatchVersion("1.14.10-gke.27")
	if e != nil || s != 10 {
		t.Fatalf("No match: %d", s)
	}
}

func TestPatchVersion2(t *testing.T) {
	s, e := PatchVersion("1.14.10")
	if e != nil || s != 10 {
		t.Fatalf("No match: %d", s)
	}
}

func TestPatchVersionErr(t *testing.T) {
	s, e := PatchVersion("1.14")
	if e != nil {
		t.Fatal(e)
	}
	assert.Equal(t, NoPatchSpecified, s)
}
