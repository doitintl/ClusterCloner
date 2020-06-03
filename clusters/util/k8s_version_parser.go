package util

import (
	"github.com/pkg/errors"
	"math"
	"regexp"
	"strconv"
)

var regexMajorMinorAndOptionalPatch = regexp.MustCompile(`^\d+\.\d+(\.\d+)?`)
var regexMajorMinorOnly = regexp.MustCompile(`^\d+\.\d+$`)
var regexMajorMinorPfx = regexp.MustCompile(`^\d+\.\d+`)
var regexMajorMinorAndDotBeforePatch = regexp.MustCompile(`^\d+\.\d+\.(\d+)`)

// MajorMinorPatchVersion ...
func MajorMinorPatchVersion(fullVersion string) (vers string, err error) {

	match := regexMajorMinorAndOptionalPatch.FindString(fullVersion)
	if match == "" {
		return "", errors.New("No match on " + fullVersion)
	}
	majorMinorOnly := regexMajorMinorOnly.FindString(fullVersion)
	if majorMinorOnly != "" {
		return match + ".0", nil
	}
	return match, nil

}

// MajorMinorVersion ...
func MajorMinorVersion(fullVersion string) (vers string, err error) {
	match := regexMajorMinorPfx.FindString(fullVersion)
	if match == "" {
		return "", errors.New("No match on " + fullVersion)
	}
	return match, nil

}

// NoPatchSpecified ...
var NoPatchSpecified = -1

// PatchVersion ...
func PatchVersion(fullVersion string) (int, error) {
	match := regexMajorMinorAndDotBeforePatch.FindStringSubmatch(fullVersion)
	if match == nil {
		match = regexMajorMinorOnly.FindStringSubmatch(fullVersion)
		if match != nil {
			return NoPatchSpecified, nil //use -1 for nil, indicating no separate patch version
		}
		return math.MinInt32, errors.New("No match on " + fullVersion)
	}
	captureGroup := match[1]
	ret, err := strconv.Atoi(captureGroup)
	if err != nil {
		panic(err) //should not happen given the regex
	}
	return ret, nil

}
