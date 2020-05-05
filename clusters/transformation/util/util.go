package util

import (
	"clusterCloner/clusters/cluster_info"
	"github.com/pkg/errors"
	"regexp"
)

func TransformSpoke(in cluster_info.ClusterInfo, outputScope, targetCloud, targetLoc, k8sVersion string) cluster_info.ClusterInfo {
	var ret = in
	ret.SourceCluster = &in
	ret.GeneratedBy = cluster_info.TRANSFORMATION
	if in.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = targetCloud
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = outputScope
	ret.Location = targetLoc
	ret.K8sVersion = in.K8sVersion
	return ret
}

func MajorMinorPatchVersion(fullVersion string) (string, error) {
	re := regexp.MustCompile(`^\d+\.\d+(\.\d+)?`)
	re2 := regexp.MustCompile(`^\d+\.\d+$`)
	match := re.FindString(fullVersion)
	if match == "" {
		return "", errors.New("No match on " + fullVersion)
	} else {
		majorMinorOnly := re2.FindString(fullVersion)
		if majorMinorOnly != "" {
			return match + ".0", nil
		}
		return match, nil

	}
}
