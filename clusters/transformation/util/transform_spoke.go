package util

import (
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/transformation/nodes"
	"github.com/google/martian/log"
)

// TransformSpoke ...
func TransformSpoke(in *clusterinfo.ClusterInfo, outputScope, targetCloud, targetLoc, k8sVersion string, machineTypes map[string]clusterinfo.MachineType) *clusterinfo.ClusterInfo {
	var ret = &clusterinfo.ClusterInfo{}
	ret.Name = in.Name
	ret.SourceCluster = in
	ret.GeneratedBy = clusterinfo.TRANSFORMATION
	ret.Cloud = targetCloud
	// ret.Name unchanged
	// ret.DeprecatedNodeCount unchanged
	ret.Scope = outputScope
	ret.Location = targetLoc
	ret.K8sVersion = in.K8sVersion
	ret.NodePools = make([]clusterinfo.NodePoolInfo, 0)
	for _, nodePoolIn := range in.NodePools {
		nodePoolOut := nodes.TransformNode(nodePoolIn, machineTypes)
		zero := clusterinfo.NodePoolInfo{}
		if nodePoolOut == zero {
			log.Errorf("Empty result of converting %v", nodePoolIn)
			return nil
		}

		ret.AddNodePool(nodePoolOut)
	}

	return ret
}
