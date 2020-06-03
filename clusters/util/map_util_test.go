package util

import (
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
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
		t.Fatalf("%s !=%s", out2Str, outS)
	}
}

func TestReverseOrderedMap(t *testing.T) {
	m := orderedmap.New()
	m.Set("a", "aa")
	m.Set("b", "bb")
	m.Set("a", "aaa")
	m.Set("c", "cc")
	m.Set("B", "bb")
	reversed := ReverseOrderedMap(m)
	s := ""
	for _, k := range reversed.Keys() {
		v, ok := reversed.Get(k)
		assert.True(t, ok)
		s += fmt.Sprintf("%v:%v,", k, v)
	}
	s=s[:len(s)-1]
	assert.Equal(t, "aaa:a,bb:B,cc:c", s)
}
