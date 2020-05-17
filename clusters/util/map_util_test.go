package util

import (
	"testing"
)

func TestStrMapToStrPtrMap(t *testing.T) {

	m := map[string]string{
		"a": "b",
		"c": "",
	}

	out := StrMapToStrPtrMap(m)
	outS := ToJSON(out)

	out2 := StrPtrMapToStrMap(out)
	out2Str := ToJSON(out2)

	if out2["a"] != m["a"] && out2["c"] != m["c"] {
		t.Errorf("%s !=%s", out2Str, outS)
	}
}
