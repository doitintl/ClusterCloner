package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/aks/access"
	transformutil "clustercloner/clusters/transformation/util"
	clusterutil "clustercloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
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

	ret := transformutil.TransformSpoke(in, "", clusters.HUB, loc, clusterK8sVersion, nil)

	return ret, err
}

// HubToCloud ...
func (tr *AKSTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting location")
	}
	ret := transformutil.TransformSpoke(in, outputScope, clusters.AZURE, loc, in.K8sVersion, access.MachineTypes)
	err = fixAksK8sVersion(ret)
	if err != nil {
		return nil, errors.Wrap(err, "error in  fixing AKS supported version")
	}
	return ret, err
}

//todo this is not a good way to fix up the node pools. In fact, we should fix K8s Version before transforming NodePools
func fixAksK8sVersion(ci *clusters.ClusterInfo) error {

	var err error
	ci.K8sVersion, err = access.FindBestMatchingSupportedK8sVersion(ci.K8sVersion)
	if err != nil {
		return errors.Wrap(err, "cannot find matching AKS version")
	}
	nodePools := ci.NodePools[:]
	ci.NodePools = make([]clusters.NodePoolInfo, 0)
	for _, np := range nodePools {
		newNp := np
		newNp.K8sVersion, err = access.FindBestMatchingSupportedK8sVersion(np.K8sVersion)
		if err != nil {
			return errors.Wrap(err, "cannot find matching AKS version")
		}
		ci.AddNodePool(newNp)
	}
	return nil

}

//LocationCloudToHub ...
func (*AKSTransformer) LocationCloudToHub(loc string) (string, error) {
	mapping, err := getAzureToHubLocations()
	if err != nil {
		return "", err
	}
	hubValue, wasinMap := mapping[loc]
	if !wasinMap {
		return "", errors.New(fmt.Sprintf("Not found: %s", loc))
	}
	return hubValue, nil
}

func getAzureToHubLocations() (map[string]string, error) {
	ret := make(map[string]string)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PWD", dir)
	fn := clusterutil.RootPath() + "/locations/azure_locations.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		wd, _ := os.Getwd()
		log.Println("At ", wd, ":", err)
		return nil, err
	}

	r := csv.NewReader(csvfile)
	r.Comma = ';'
	first := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if first {
			first = false
			continue
		}
		if len(record) == 1 {
			log.Println("Short record ", record)
		}
		azRegion := record[3]
		hubRegion := record[5]
		supportsAks := record[4]
		if supportsAks != "true" {
			return nil, errors.New(fmt.Sprintf("Azure region %s does not support AKS", azRegion))
		}
		ret[azRegion] = hubRegion
	}
	return ret, nil
}

//LocationHubToCloud ...
func (AKSTransformer) LocationHubToCloud(location string) (string, error) {
	azToHub, err := getAzureToHubLocations()
	if err != nil {
		return "", err
	}
	hubToAz := reverseMap(azToHub)
	azLoc, ok := hubToAz[location]
	if !ok {
		return "", errors.New(fmt.Sprintf("Cannot find %s", location))
	}
	return azLoc, nil

}

func reverseMap(m map[string]string) map[string]string {
	reverse := make(map[string]string)
	var dupes = make([][3]string, 0)
	for k, v := range m {
		existing, wasInMap := reverse[v]
		if wasInMap {
			var using, notUsing string
			if k < existing {
				using = k
				notUsing = existing
			} else {
				using = existing
				notUsing = k
			}
			dupeTriple := [3]string{v, using, notUsing}
			dupes = append(dupes, dupeTriple)

			reverse[v] = using
		} else {
			reverse[v] = k
		}
	}
	dupesStr := ""
	for _, triple := range dupes {
		dupesStr += "New Key \"" + triple[0] + "\"; key as new value \"" + triple[1] + "\"; Key not used as new value \"" + triple[2] + "\"\n"
	}
	log.Println("Duplicates in reversing map: ", dupesStr)
	return reverse
}
