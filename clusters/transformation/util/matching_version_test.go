package util

import "testing"

func TestSupportedK8sVersion(t *testing.T) {
	supportedVersions := []string{"1.14.8", "1.14.9", "1.14.11", "1.15.1"}
	matchingSupported, err := findBestMatchingSupportedK8sVersion("1.14.1", supportedVersions)
	if err != nil {
		t.Error(err)
	}
	if matchingSupported != "1.14.8" {
		t.Error(matchingSupported)
	}
}

func TestSupportedK8sVersionError(t *testing.T) {
	supportedVersions := []string{"1.14.8", "1.14.9", "1.14.11", "1.15.1"}

	_, err := findBestMatchingSupportedK8sVersion("1.214.10", supportedVersions)
	if err == nil {
		t.Error(err)
	}

}

func TestSupportedK8sVersion3(t *testing.T) {
	supportedVersions := []string{"1.14.8", "1.14.9", "1.14.11", "1.15.1"}

	supported, err := findBestMatchingSupportedK8sVersion("1.14.10", supportedVersions)
	if err != nil {
		t.Error(err)
	}
	if supported != "1.14.11" {
		t.Error(supported)
	}
}
