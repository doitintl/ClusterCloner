package transform

import (
	"clustercloner/clusters"
)

// EKSTransformer ...
type EKSTransformer struct{}

// CloudToHub ...
func (tr *EKSTransformer) CloudToHub(in *clusters.ClusterInfo) (ret *clusters.ClusterInfo, err error) {

	return ret, err
}

// HubToCloud ...
func (tr *EKSTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (ret *clusters.ClusterInfo, err error) {
	return ret, err
}

//LocationCloudToHub ...
func (*EKSTransformer) LocationCloudToHub(loc string) (hubValue string, err error) {
	return hubValue, nil
}

//LocationHubToCloud ...
func (EKSTransformer) LocationHubToCloud(location string) (ret string, err error) {

	return ret, nil

}
