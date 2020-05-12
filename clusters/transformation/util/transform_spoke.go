package util

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess"
	"clustercloner/clusters/transformation/nodes"
	clusterutil "clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math"
)

// TransformSpoke ...
func TransformSpoke(in *clusters.ClusterInfo, outputScope, targetCloud, targetLoc,
	targetClusterK8sVersion string, machineTypes map[string]clusters.MachineType,
	adjustK8sVersions bool) *clusters.ClusterInfo {

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
	if adjustK8sVersions {
		err := fixK8sVersion(ret) //should not fix version post-facto like this
		if err != nil {
			log.Println(err, "cannot fix K8s versions")
			return nil
		}
	}
	return ret

}

// fixK8sVersion ...
func fixK8sVersion(ci *clusters.ClusterInfo) error {
	ca := clusteraccess.GetClusterAccessor(ci.Cloud)
	if ca == nil {
		return errors.New("cannot get cluster accessor for " + ci.Cloud)
	}
	supportedVersions := ca.GetSupportedK8sVersions(ci.Scope, ci.Location)
	if supportedVersions == nil {
		return errors.New("cannot find supported K8s versions")
	}
	var err error
	ci.K8sVersion, err = findBestMatchingSupportedK8sVersion(ci.K8sVersion, supportedVersions)
	if err != nil {
		return errors.Wrap(err, "cannot find matching AKS version")
	}
	nodePools := ci.NodePools[:]
	ci.NodePools = make([]clusters.NodePoolInfo, 0)
	for _, np := range nodePools {
		newNp := np
		newNp.K8sVersion, err = findBestMatchingSupportedK8sVersion(np.K8sVersion, supportedVersions)
		if err != nil {
			return errors.Wrap(err, "cannot find matching AKS version")
		}
		ci.AddNodePool(newNp)
	}
	return nil

}

/*FindBestMatchingSupportedK8sVersion  find the least patch version that is
greater or equal to  the supplied vers, but has the same major-minor version.
If that not possible, get the largest patch version that has the same major-minor version
*/
func findBestMatchingSupportedK8sVersion(vers string, supportedVersions []string) (string, error) {
	var potentialMatchPatchVersion = math.MaxInt32
	majorMinor, err := clusterutil.MajorMinorVersion(vers)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse versions")
	}
	patchV, err := clusterutil.PatchVersion(vers)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse versions")
	}
	for _, supported := range supportedVersions {
		majorMinorSupported, err := clusterutil.MajorMinorVersion(supported)
		if err != nil {
			return "", errors.Wrap(err, "cannot parse versions")
		}
		if majorMinor == majorMinorSupported {
			var patchSupported int
			patchSupported, err = clusterutil.PatchVersion(supported)
			if err != nil {
				panic(err) //should not happen
			}
			if patchSupported < potentialMatchPatchVersion && patchSupported >= patchV {
				potentialMatchPatchVersion = patchSupported
			}
		}
	}
	if potentialMatchPatchVersion == math.MaxInt32 {
		potentialMatchPatchVersion = math.MinInt32
		//get largest patch version in this major-minor
		for _, supported := range supportedVersions {
			majorMinorSupported, err := clusterutil.MajorMinorVersion(supported)
			if err != nil {
				return "", errors.Wrap(err, "cannot parse versions")
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
			return "", errors.New("cannot match to patch version: " + vers)

		}
	}
	ret := fmt.Sprintf("%s.%d", majorMinor, potentialMatchPatchVersion)
	return ret, nil
}
