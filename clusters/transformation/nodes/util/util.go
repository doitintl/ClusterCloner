package util

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

// MajorMinorPatchVersion ...
func MajorMinorPatchVersion(fullVersion string) (vers string, err error) {
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

// MajorMinorVersion ...
func MajorMinorVersion(fullVersion string) (vers string, err error) {
	re := regexp.MustCompile(`^\d+\.\d+`)
	match := re.FindString(fullVersion)
	if match == "" {
		return "", errors.New("No match on " + fullVersion)
	}
	return match, nil

}

// PatchVersion ...
func PatchVersion(fullVersion string) (int, error) {
	re := regexp.MustCompile(`^\d+\.\d+\.(\d+)?`)
	match := re.FindStringSubmatch(fullVersion)
	if match == nil {
		return -1, errors.New("No match on " + fullVersion)
	}
	captureGroup := match[1]
	ret, err := strconv.Atoi(captureGroup)
	if err != nil {
		panic(err) //should not happen given the regex
	}
	return ret, nil

}
