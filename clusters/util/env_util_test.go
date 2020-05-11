package util

import (
	"strings"
	"testing"
)

func TestRootPath(t *testing.T) {

	r := RootPath()
	if !strings.Contains(r, "clustercloner") {
		t.Error(r)
	}
}
