package util

import (
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/transformation/nodes"
)

// TransformSpoke ...
func TransformSpoke(in clusterinfo.ClusterInfo, outputScope, targetCloud, targetLoc, k8sVersion string, machineTypes map[string]clusterinfo.MachineType) clusterinfo.ClusterInfo {
	var ret = in
	ret.SourceCluster = &in
	ret.GeneratedBy = clusterinfo.TRANSFORMATION
	if in.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = targetCloud
	// ret.Name unchanged
	// ret.DeprecatedNodeCount unchanged
	ret.Scope = outputScope
	ret.Location = targetLoc
	ret.K8sVersion = in.K8sVersion
	for _, np := range in.NodePools {
		ret.AddNodePool(nodes.TransformNode(np, machineTypes))
	}
	return ret
}
