package nodes

import (
	"clustercloner/clusters"
	"clustercloner/clusters/util"
	"github.com/pkg/errors"
	"math"
)

// MachineTypeConverter ...
type MachineTypeConverter func(mt clusters.MachineType) clusters.MachineType

// TransformNodePool ...
func TransformNodePool(np clusters.NodePoolInfo, machineTypes map[string]clusters.MachineType) (clusters.NodePoolInfo, error) {
	nodePoolK8sVersion, err := util.MajorMinorPatchVersion(np.K8sVersion)
	if err != nil {
		return clusters.NodePoolInfo{}, errors.New("cannot convert K8s Version \"" + np.K8sVersion + "\" for node pool")
	}
	ret := clusters.NodePoolInfo{
		Name:        np.Name,
		NodeCount:   np.NodeCount,
		K8sVersion:  nodePoolK8sVersion,
		MachineType: FindMatchingMachineType(np.MachineType, machineTypes),
		DiskSizeGB:  np.DiskSizeGB,
	}
	return ret, nil
}

// FindMatchingMachineType chooses the weakest machine wgich equals or exceeeds in the inputMachineType in CPU amd RAM. If there are several some such, one is chosen arbitrarily
func FindMatchingMachineType(inputMachineType clusters.MachineType, machineTypes map[string]clusters.MachineType) clusters.MachineType {
	if machineTypes == nil { //Transforming to Hub, no change
		return inputMachineType
	}
	leastUpperBound := clusters.MachineType{Name: "<NONE KNOWN>", CPU: math.MaxInt32, RAMGB: math.MaxInt32}
	for _, candidateMachineType := range machineTypes {
		if candidateMachineType.RAMGB >= inputMachineType.RAMGB && candidateMachineType.CPU >= inputMachineType.CPU {
			if candidateMachineType.RAMGB <= leastUpperBound.RAMGB && candidateMachineType.CPU <= leastUpperBound.CPU {
				leastUpperBound = candidateMachineType
			}
		}
	}
	return leastUpperBound
}
