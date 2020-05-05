package transform

import (
	"clusterCloner/clusters/cluster_info"
	transformutil "clusterCloner/clusters/transformation/util"
	clusterutil "clusterCloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type AksTransformer struct{}

func (tr AksTransformer) CloudToHub(in cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	loc, err := tr.LocationCloudToHub(in.Location)
	if err != nil {
		return cluster_info.ClusterInfo{}, errors.Wrap(err, "error in converting locations")
	}

	k8sVersion, err := transformutil.MajorMinorPatchVersion(in.K8sVersion)
	if err != nil {
		return cluster_info.ClusterInfo{}, errors.Wrap(err, "error in K8s Version "+in.K8sVersion)
	}

	ret := transformutil.TransformSpoke(in, "", cluster_info.HUB, loc, k8sVersion)

	return ret, err
}

func (tr AksTransformer) HubToCloud(in cluster_info.ClusterInfo, outputScope string) (cluster_info.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return cluster_info.ClusterInfo{}, errors.Wrap(err, "error in converting location")
	}
	ret := transformutil.TransformSpoke(in, outputScope, cluster_info.AZURE, loc, in.K8sVersion)
	ret.Name = ret.Name + "arbitrarysuffix"
	return ret, err
}

func (AksTransformer) LocationCloudToHub(loc string) (string, error) {
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
			log.Print("Short record ", record)
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
func (AksTransformer) LocationHubToCloud(location string) (string, error) {
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

	log.Println("Duplicates in reversing map: New keys, followed by a value (old key) that will be used and one that won't ", dupes)
	return reverse
}
