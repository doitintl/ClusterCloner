package util

import (
	"clustercloner/clusters"
	"clustercloner/clusters/util"
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

// CopyClusterInfo ...
func CopyClusterInfo(in *clusters.ClusterInfo) clusters.ClusterInfo {
	ret := *in //copy
	ret.Labels = util.CopyStringMap(in.Labels)
	ret.NodePools = copyNodePools(in) //array isa pointer type, needs copying
	return ret
}

// CloudToHub ...
func (it *IdentityTransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	ret := CopyClusterInfo(in)
	ret.SourceCluster = in
	ret.Cloud = clusters.Hub
	ret.GeneratedBy = clusters.Transformation
	return &ret, nil
}

// HubToCloud ...
func (it *IdentityTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	ret := CopyClusterInfo(in)
	ret.SourceCluster = in
	ret.GeneratedBy = clusters.Transformation
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
