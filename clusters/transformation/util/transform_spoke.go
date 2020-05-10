package util

import (
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/transformation/nodes"
	"log"
)

// TransformSpoke ...
func TransformSpoke(in *clusterinfo.ClusterInfo, outputScope, targetCloud, targetLoc, targetClusterK8sVersion string, machineTypes map[string]clusterinfo.MachineType) *clusterinfo.ClusterInfo {
	var ret = &clusterinfo.ClusterInfo{}
	ret.Name = in.Name
	ret.SourceCluster = in
	ret.GeneratedBy = clusterinfo.TRANSFORMATION
	ret.Cloud = targetCloud
	// ret.Name unchanged
	// ret.DeprecatedNodeCount unchanged
	ret.Scope = outputScope
	ret.Location = targetLoc
	ret.K8sVersion = targetClusterK8sVersion
	ret.NodePools = make([]clusterinfo.NodePoolInfo, 0)
	for _, nodePoolIn := range in.NodePools {
		nodePoolOut := nodes.TransformNodePool(nodePoolIn, machineTypes)
		zero := clusterinfo.NodePoolInfo{}
		if nodePoolOut == zero {
			log.Printf("Empty result of converting %v", nodePoolIn)
			return nil
		}

		ret.AddNodePool(nodePoolOut)
	}

	return ret
}
