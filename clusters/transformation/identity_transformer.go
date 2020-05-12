package transformation

import (
	"clustercloner/clusters"
	"github.com/pkg/errors"
)

// IdentityTransformer ...
type IdentityTransformer struct{ TargetCloud string }

func copyNodePools(in *clusters.ClusterInfo) []clusters.NodePoolInfo {
	copyNPs := make([]clusters.NodePoolInfo, len(in.NodePools))
	for i := range in.NodePools {
		copyNPs[i] = in.NodePools[i] //copy value
	}
	return copyNPs
}

func copyClusterInfo(in *clusters.ClusterInfo) clusters.ClusterInfo {
	retVal := *in                        //copy
	retVal.NodePools = copyNodePools(in) //array isa pointer type, needs copying
	return retVal
}

// CloudToHub ...
func (it *IdentityTransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	ret := copyClusterInfo(in)
	ret.Cloud = clusters.HUB
	ret.GeneratedBy = clusters.TRANSFORMATION
	return &ret, nil
}

// HubToCloud ...
func (it *IdentityTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	ret := copyClusterInfo(in)
	ret.GeneratedBy = clusters.TRANSFORMATION
	if it.TargetCloud == "" {
		return nil, errors.New("No TargetCloud specified in IdentityTransformer")
	}
	ret.Cloud = it.TargetCloud
	ret.Scope = outputScope
	return &ret, nil
}

// LocationHubToCloud ...
func (it *IdentityTransformer) LocationHubToCloud(loc string) (string, error) {
	return loc, nil
}

// LocationCloudToHub ...
func (it *IdentityTransformer) LocationCloudToHub(loc string) (string, error) {
	return loc, nil
}
