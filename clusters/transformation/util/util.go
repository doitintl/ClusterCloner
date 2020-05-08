package util

import (
	"github.com/pkg/errors"
	"regexp"
)

// MajorMinorPatchVersion ...
func MajorMinorPatchVersion(fullVersion string) (string, error) {
	re := regexp.MustCompile(`^\d+\.\d+(\.\d+)?`)
	re2 := regexp.MustCompile(`^\d+\.\d+$`)
	match := re.FindString(fullVersion)
	if match == "" {
		return "", errors.New("No match on " + fullVersion)
	}
	majorMinorOnly := re2.FindString(fullVersion)
	if majorMinorOnly != "" {
		return match + ".0", nil
	}
	return match, nil

}
