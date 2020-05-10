package nodes

import (
	"clustercloner/clusters/clusterinfo"
	nodeutil "clustercloner/clusters/transformation/nodes/util"
	"math"
)

// MachineTypeConverter ...
type MachineTypeConverter func(mt clusterinfo.MachineType) clusterinfo.MachineType

// TransformNodePool ...
func TransformNodePool(np clusterinfo.NodePoolInfo, machineTypes map[string]clusterinfo.MachineType) clusterinfo.NodePoolInfo {
	nodePoolK8sVersion, err := nodeutil.MajorMinorPatchVersion(np.K8sVersion)
	if err != nil {
		panic(err) //todo fix
	}
	ret := clusterinfo.NodePoolInfo{
		Name:        np.Name,
		NodeCount:   np.NodeCount,
		K8sVersion:  nodePoolK8sVersion,
		MachineType: FindMatchingMachineType(np.MachineType, machineTypes),
		DiskSizeGB:  np.DiskSizeGB,
	}
	return ret
}

// FindMatchingMachineType chooses the weakest machine wgich equals or exceeeds in the inputMachineType in CPU amd RAM. If there are several some such, one is chosen arbitrarily
func FindMatchingMachineType(inputMachineType clusterinfo.MachineType, machineTypes map[string]clusterinfo.MachineType) clusterinfo.MachineType {
	if machineTypes == nil { //Transforming to Hub, no change
		return inputMachineType
	}
	leastUpperBound := clusterinfo.MachineType{Name: "<NONE KNOWN>", CPU: math.MaxInt32, RAMGB: math.MaxInt32}
	for _, candidateMachineType := range machineTypes {
		if candidateMachineType.RAMGB >= inputMachineType.RAMGB && candidateMachineType.CPU >= inputMachineType.CPU {
			if candidateMachineType.RAMGB <= leastUpperBound.RAMGB && candidateMachineType.CPU <= leastUpperBound.CPU {
				leastUpperBound = candidateMachineType
			}
		}
	}
	return leastUpperBound
}
