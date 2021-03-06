package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/aks/access"
	transformutil "clustercloner/clusters/transformation/util"
	clusterutil "clustercloner/clusters/util"
	"github.com/iancoleman/orderedmap"
	"github.com/pkg/errors"
)

// AKSTransformer ...
type AKSTransformer struct{}

// CloudToHub ...
func (tr *AKSTransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
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
		return nil, errors.Wrap(err, "cannot TransformSpoke CloudToHub AKS")
	}
	return ret, nil
}

// HubToCloud ...
func (tr *AKSTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting location")
	}
	ret, err := transformutil.TransformSpoke(in, outputScope, clusters.Azure, loc, in.K8sVersion, access.GetMachineTypes(), true)

	if err != nil {
		return nil, errors.Wrap(err, "cannot TransformSpoke HubToCloud AKS")
	}
	return ret, nil
}

//LocationCloudToHub ...
func (*AKSTransformer) LocationCloudToHub(loc string) (string, error) {
	mapping, err := locationsCloudToHub()
	if err != nil {
		return "", errors.Wrap(err, "error getting locationsCloudToHub")
	}
	hubValue, wasinMap := mapping.Get(loc)
	if !wasinMap {
		return "", errors.Errorf("Not found: %s", loc)
	}
	hubValStr, ok := hubValue.(string)
	if !ok {
		panic(hubValue)
	}
	return hubValStr, nil
}

var locations *orderedmap.OrderedMap

// locationsCloudToHub ...
func locationsCloudToHub() (orderedMap *orderedmap.OrderedMap, string error) {
	file := "azure_locations.csv"
	if locations == nil {
		var err error
		locations, err = transformutil.LoadLocationMap(file)
		if err != nil {
			return nil, errors.Wrap(err, "cannot load "+file)
		}
	}
	return locations, nil
}

//LocationHubToCloud ...
func (AKSTransformer) LocationHubToCloud(location string) (string, error) {
	azToHub, err := locationsCloudToHub()
	if err != nil {
		return "", errors.Wrap(err, "cannot get locationsCloudToHub Azure")

	}
	hubToAz := clusterutil.ReverseOrderedMap(azToHub)
	azLoc, ok := hubToAz.Get(location)
	if !ok {
		return "", errors.Errorf("Cannot find %s", location)
	}
	return azLoc.(string), nil

}
