package util

// Keys ...
func Keys(m map[string]string) []string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ContainsInt ...
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsStr ...
func ContainsStr(slice []string, elem string) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}

// Contains ...
func Contains(slice []interface{}, elem interface{}) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}
