package util

import (
	"clustercloner/clusters"
	"clustercloner/clusters/transformation/nodes"
	"log"
)

// TransformSpoke ...
func TransformSpoke(in *clusters.ClusterInfo, outputScope, targetCloud, targetLoc, targetClusterK8sVersion string, machineTypes map[string]clusters.MachineType) *clusters.ClusterInfo {
	var ret = &clusters.ClusterInfo{}
	ret.Name = in.Name
	ret.SourceCluster = in
	ret.GeneratedBy = clusters.TRANSFORMATION
	ret.Cloud = targetCloud
	// ret.Name unchanged
	// ret.DeprecatedNodeCount unchanged
	ret.Scope = outputScope
	ret.Location = targetLoc
	ret.K8sVersion = targetClusterK8sVersion
	ret.NodePools = make([]clusters.NodePoolInfo, 0)
	for _, nodePoolIn := range in.NodePools {
		nodePoolOut, err := nodes.TransformNodePool(nodePoolIn, machineTypes)
		if err != nil {
			log.Printf("Error transforming Node Pool %v\n", err)
			return nil
		}
		zero := clusters.NodePoolInfo{}
		if nodePoolOut == zero {
			log.Printf("Empty result of converting %v", nodePoolIn)
			return nil
		}

		ret.AddNodePool(nodePoolOut)
	}

	return ret
}
