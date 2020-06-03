package util

import (
	"fmt"
	"github.com/iancoleman/orderedmap"
	"log"
)

// LabelMatch ...
func LabelMatch(labelFilter map[string]string, actualLabels map[string]string) bool {

	for k, v := range labelFilter {
		actualVal, found := actualLabels[k]
		if !(found && actualVal == v) {
			return false
		}
	}

	return true
}

// CopyStringMap ...
func CopyStringMap(m map[string]string) map[string]string {
	cp := make(map[string]string)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}

// StrMapToStrPtrMap ...
func StrMapToStrPtrMap(m map[string]string) map[string]*string {
	ret := make(map[string]*string)

	for k, v := range m {
		v2 := v
		ret[k] = &v2
	}
	return ret
}

// StrPtrMapToStrMap ...
func StrPtrMapToStrMap(m map[string]*string) map[string]string {
	ret := make(map[string]string)
	for k, v := range m {
		ret[k] = *v
	}
	return ret

}

// ReverseOrderedMap ...
func ReverseOrderedMap(m *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	reverse := orderedmap.New()
	for _, k := range m.Keys() {
		v, ok := m.Get(k)
		if !ok {
			panic(k)
		}
		vStr, wasStr := v.(string)
		if !wasStr {
			panic(fmt.Sprintf("expect string values %v", v))
		}
		reverse.Set(vStr, k)
	}
	return reverse
}

// ReverseStrMap ...
func ReverseStrMap(m map[string]string) map[string]string {
	reverse := make(map[string]string)
	var dupes = make([][3]string, 0)
	for k, v := range m {
		existing, wasInMap := reverse[v]
		if wasInMap {
			var using, notUsing string
			if k < existing {
				using = k
				notUsing = existing
			} else {
				using = existing
				notUsing = k
			}
			dupeTriple := [3]string{v, using, notUsing}
			dupes = append(dupes, dupeTriple)

			reverse[v] = using
		} else {
			reverse[v] = k
		}
	}
	dupesStr := ""
	for _, triple := range dupes {
		dupesStr += "Key \"" + triple[0] + "\"; using value \"" + triple[1] + "\"; dropping value \"" + triple[2] + "\"; "
	}
	log.Println("Duplicates in reversing map: ", dupesStr)
	return reverse
}
