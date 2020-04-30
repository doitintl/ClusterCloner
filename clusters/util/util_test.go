package util

import (
	"strings"
	"testing"
)

func TestRootPath(t *testing.T) {
	r := RootPath()
	_ = r
	if !strings.Contains(r, "clustercloner") {
		t.Error(r)
	}

}
