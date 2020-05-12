package transformation

import "clustercloner/clusters"

// IdentityTransformer ...
type IdentityTransformer struct{}

// CloudToHub ...
func (it *IdentityTransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	copyNPs := make([]clusters.NodePoolInfo, len(in.NodePools))
	for i := range in.NodePools {
		copyNPs[i] = in.NodePools[i] //copy value
	}
	ret := clusters.ClusterInfo{
		Cloud:         clusters.HUB,
		Scope:         in.Scope,
		Location:      in.Location,
		Name:          in.Name,
		K8sVersion:    in.K8sVersion,
		GeneratedBy:   clusters.TRANSFORMATION,
		NodePools:     copyNPs,
		SourceCluster: in,
	}
	return &ret, nil
}

// HubToCloud ...
func (it *IdentityTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	ret := clusters.ClusterInfo{
		Cloud:         clusters.HUB,
		Scope:         outputScope,
		Location:      in.Location,
		Name:          in.Name,
		K8sVersion:    in.K8sVersion,
		GeneratedBy:   clusters.TRANSFORMATION,
		NodePools:     in.NodePools[:],
		SourceCluster: in,
	}
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
