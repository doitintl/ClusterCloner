package util

import (
	"reflect"
	"testing"
)

func TestMapToCsv(t *testing.T) {

	s := "a=b,cc=ddd,eee=f,g,h="
	m := CommaSeparatedKeyValPairsToMap(s)
	s2 := ToCommaSeparateKeyValuePairs(m)
	m2 := CommaSeparatedKeyValPairsToMap(s2)
	if !reflect.DeepEqual(m, m2) {
		t.Fatal(m2)
	}

}
