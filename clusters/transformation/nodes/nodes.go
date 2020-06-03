package nodes

import (
	"clustercloner/clusters"
	"clustercloner/clusters/machinetypes"
	"clustercloner/clusters/util"
	"github.com/pkg/errors"
	"math"
)

// MachineTypeConverter ...
type MachineTypeConverter func(mt machinetypes.MachineType) machinetypes.MachineType

// TransformNodePool ...
func TransformNodePool(np clusters.NodePoolInfo, machineTypes *machinetypes.MachineTypes) (clusters.NodePoolInfo, error) {
	nodePoolK8sVersion, err := util.MajorMinorPatchVersion(np.K8sVersion)
	if err != nil {
		return clusters.NodePoolInfo{}, errors.New("cannot convert K8s Version \"" + np.K8sVersion + "\" for node pool")
	}
	matchingMachineType := FindMatchingMachineType(np.MachineType, machineTypes)
	if matchingMachineType.Name == "" { //zero-object -- name n
		return clusters.NodePoolInfo{}, errors.New("cannot find match for Machine Type " + np.MachineType.Name)
	}
	ret := clusters.NodePoolInfo{
		Name:        np.Name,
		NodeCount:   np.NodeCount,
		K8sVersion:  nodePoolK8sVersion,
		MachineType: matchingMachineType,
		DiskSizeGB:  np.DiskSizeGB,
		Preemptible: np.Preemptible,
	}
	return ret, nil
}

// FindMatchingMachineType chooses the weakest machine which equals or exceeeds in the inputMachineType in CPU amd RAM. If there are several some such, one is chosen arbitrarily
func FindMatchingMachineType(inputMachineType machinetypes.MachineType, machineTypes *machinetypes.MachineTypes) machinetypes.MachineType {
	if machineTypes == nil { //Transforming to Hub, no change
		return inputMachineType
	}
	leastUpperBound := machinetypes.MachineType{Name: "<NONE KNOWN>", CPU: math.MaxInt32, RAMMB: math.MaxInt32}
	for _, candidateMachineType := range machineTypes.List() {
		if candidateMachineType.RAMMB >= inputMachineType.RAMMB && candidateMachineType.CPU >= inputMachineType.CPU {
			if candidateMachineType.RAMMB <= leastUpperBound.RAMMB && candidateMachineType.CPU <= leastUpperBound.CPU {
				leastUpperBound = candidateMachineType
			}
		}
	}
	return leastUpperBound
}
