package util

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
