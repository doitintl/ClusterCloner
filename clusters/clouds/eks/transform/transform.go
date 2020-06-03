package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/access"
	transformutil "clustercloner/clusters/transformation/util"
	clusterutil "clustercloner/clusters/util"
	"github.com/pkg/errors"
)

// EKSTransformer ...
type EKSTransformer struct{}

// CloudToHub ...
func (tr *EKSTransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationCloudToHub(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting locations")
	}

	clusterK8sVersion, err := clusterutil.MajorMinorPatchVersion(in.K8sVersion)
	if err != nil {
		return nil, errors.Wrap(err, "error in K8s K8sVersion "+in.K8sVersion)
	}

	ret, err := transformutil.TransformSpoke(in, "", clusters.Hub, loc, clusterK8sVersion, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "cannot TransformSpoke CloudToHub EKS")
	}
	return ret, nil
}

// HubToCloud ...
func (tr *EKSTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting location")
	}
	ret, err := transformutil.TransformSpoke(in, outputScope, clusters.AWS, loc, in.K8sVersion, access.GetMachineTypes(), true)

	if err != nil {
		return nil, errors.Wrap(err, "cannot TransformSpoke HubToCloud EKS")
	}
	return ret, nil
}

// LocationsCloudToHub ...
func LocationsCloudToHub() (map[string]string, error) {
	file := "aws_locations.csv"
	if locations == nil {
		var err error
		locations, err = transformutil.LoadLocationMap(file)
		if err != nil {
			return nil, errors.Wrap(err, "cannot load "+file)
		}
	}
	return locations, nil
}

//LocationCloudToHub ...
func (*EKSTransformer) LocationCloudToHub(loc string) (string, error) {
	mapping, err := LocationsCloudToHub()
	if err != nil {
		return "", errors.Wrap(err, "error getting LocationsCloudToHub")
	}
	hubValue, wasinMap := mapping[loc]
	if !wasinMap {
		return "", errors.Errorf("Not found: %s", loc)
	}
	return hubValue, nil
}

//LocationHubToCloud ...
func (EKSTransformer) LocationHubToCloud(location string) (string, error) {
	awsToHub, err := LocationsCloudToHub()
	if err != nil {
		return "", errors.Wrap(err, "cannot get LocationsCloudToHub AWS")
	}
	hubToAws := clusterutil.ReverseStrMap(awsToHub) // //TODO make it deterministic
	azLoc, ok := hubToAws[location]
	if !ok {
		return "", errors.Errorf("Cannot find %s", location)
	}
	return azLoc, nil

}

// Locations ...
var locations map[string]string
