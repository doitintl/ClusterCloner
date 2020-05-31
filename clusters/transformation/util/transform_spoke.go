package util

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation/nodes"
	clusterutil "clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"
	"math"
)

// TransformSpoke ...
func TransformSpoke(in *clusters.ClusterInfo, outputScope, targetCloud, targetLoc,
	clusterK8sVersion string, machineTypes map[string]clusters.MachineType,
	adjustK8sVersions bool) (*clusters.ClusterInfo, error) {

	var ret = &clusters.ClusterInfo{
		Name:          in.Name,
		SourceCluster: in,
		GeneratedBy:   clusters.Transformation,
		Cloud:         targetCloud,
		Scope:         outputScope,
		Location:      targetLoc,
		K8sVersion:    clusterK8sVersion, //temp,replaced below
		Labels:        clusterutil.CopyStringMap(in.Labels),
	}

	ret.NodePools = make([]clusters.NodePoolInfo, 0)
	for _, nodePoolIn := range in.NodePools {
		if nodePoolIn.MachineType.Name == "" {
			return nil, errors.New("node pool " + nodePoolIn.Name + " has an uninitialized Machine Type")
		}
		nodePoolOut, err := nodes.TransformNodePool(nodePoolIn, machineTypes)
		if err != nil {
			return nil, errors.Wrap(err, "error transforming Node Pool"+nodePoolIn.Name)
		}
		zero := clusters.NodePoolInfo{}
		if nodePoolOut == zero {
			return nil, errors.New(fmt.Sprintf("Empty result of converting %v", nodePoolIn))
		}

		ret.AddNodePool(nodePoolOut)
	}
	if adjustK8sVersions {
		err := fixK8sVersion(ret) //Improve design: Should not fix version after setting it wrongly,  like this
		if err != nil {
			return nil, errors.Wrap(err, "cannot fix K8s versions")
		}
	}
	return ret, nil

}

// fixK8sVersion ...
func fixK8sVersion(mutate *clusters.ClusterInfo) error {
	ca := clusteraccess.GetClusterAccess(mutate.Cloud)
	if ca == nil {
		return errors.New("cannot get cluster access for " + mutate.Cloud)
	}
	supportedVersions, err := ca.GetSupportedK8sVersions(mutate.Scope, mutate.Location)
	if err != nil {
		return errors.Wrap(err, "cannot get SupportedK8sVersions")
	}
	mutate.K8sVersion, err = findBestMatchingSupportedK8sVersion(mutate.K8sVersion, supportedVersions)
	if err != nil {
		return errors.Wrap(err, "cannot find matching K8s version")
	}
	nodePools := mutate.NodePools[:]
	mutate.NodePools = make([]clusters.NodePoolInfo, 0)
	for _, np := range nodePools {
		newNp := np
		newNp.K8sVersion, err = findBestMatchingSupportedK8sVersion(np.K8sVersion, supportedVersions)
		if err != nil {
			return errors.Wrap(err, "cannot find matching K8s version")
		}
		mutate.AddNodePool(newNp)
	}
	return nil

}

/*FindBestMatchingSupportedK8sVersion  find the least patch version that is
greater or equal to  the supplied vers, but has the same major-minor version.
If that not possible, get the largest patch version that has the same major-minor version
*/
func findBestMatchingSupportedK8sVersion(vers string, supportedVersions []string) (bestVersion string, err error) {
	potentialMatchPatchVersion, err := leastPatchGreaterThanThisWithSameMajorMinor(vers, supportedVersions)
	if err != nil {
		return "", errors.Wrap(err, "error in finding match")
	}
	majorMinor, err := clusterutil.MajorMinorVersion(vers)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse versions")
	}
	if potentialMatchPatchVersion == math.MaxInt32 {
		potentialMatchPatchVersion, err = biggestPatchVersionSameMajorMinor(vers, supportedVersions)
		if err != nil {
			return "", errors.Wrap(err, "error in finding match")
		}
	}

	if potentialMatchPatchVersion == clusterutil.NoPatchSpecified {
		bestVersion = majorMinor
	} else {
		bestVersion = fmt.Sprintf("%s.%d", majorMinor, potentialMatchPatchVersion)
	}
	return bestVersion, nil
}

func leastPatchGreaterThanThisWithSameMajorMinor(vers string, supportedVersions []string) (int, error) {
	var potentialMatchPatchVersion = math.MaxInt32
	majorMinor, err := clusterutil.MajorMinorVersion(vers)
	if err != nil {
		return 0, errors.Wrap(err, "cannot parse versions")
	}
	patchV, err := clusterutil.PatchVersion(vers)
	if err != nil {
		return 0, errors.Wrap(err, "cannot parse versions")
	}
	for _, supported := range supportedVersions {
		majorMinorSupported, err := clusterutil.MajorMinorVersion(supported)
		if err != nil {
			return 0, errors.Wrap(err, "cannot parse versions")
		}

		if majorMinor == majorMinorSupported {
			var patchSupported int
			patchSupported, err = clusterutil.PatchVersion(supported)
			if err != nil {
				panic(err) //should not happen
			}
			if patchSupported == clusterutil.NoPatchSpecified { //as with EKS
				potentialMatchPatchVersion = clusterutil.NoPatchSpecified
			} else if patchSupported < potentialMatchPatchVersion && patchSupported >= patchV {
				potentialMatchPatchVersion = patchSupported
			}
		}
	}
	return potentialMatchPatchVersion, nil
}

func biggestPatchVersionSameMajorMinor(vers string, supportedVersions []string) (int, error) {
	majorMinor, err := clusterutil.MajorMinorVersion(vers)
	if err != nil {
		return -1, errors.Wrap(err, "cannot parse versions")
	}
	patchV, err := clusterutil.PatchVersion(vers)
	if err != nil {
		return -1, errors.Wrap(err, "cannot parse versions")
	}
	potentialMatchPatchVersion := math.MinInt32
	//get largest patch version in this major-minor
	for _, supported := range supportedVersions {
		majorMinorSupported, err := clusterutil.MajorMinorVersion(supported)
		if err != nil {
			return 0, errors.Wrap(err, "cannot parse versions")
		}
		if majorMinor == majorMinorSupported {
			var patchSupported int
			patchSupported, err = clusterutil.PatchVersion(supported)
			if err != nil {
				panic(err) //should not happen
			}
			if patchSupported > potentialMatchPatchVersion {
				if patchSupported >= patchV {
					panic(fmt.Sprintf("In this part of the search, we have already found"+
						" no supported patch versions greater than"+
						" the current patch version %d", patchSupported))
				}
				potentialMatchPatchVersion = patchSupported
			}
		}
	}
	if potentialMatchPatchVersion == math.MaxInt32 || potentialMatchPatchVersion == math.MinInt32 {
		return 0, errors.New("cannot match to patch version: " + vers)

	}
	return potentialMatchPatchVersion, nil
}
