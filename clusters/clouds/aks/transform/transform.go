package transform

import (
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type AksTransformer struct{}

func (tr AksTransformer) CloudToHub(inputClusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	var ret = inputClusterInfo
	ret.SourceCluster = &inputClusterInfo
	ret.GeneratedBy = cluster_info.TRANSFORMATION
	if inputClusterInfo.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.HUB
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in conversion cross-cloud
	loc, err := tr.LocationCloudToHub(ret.Location)
	if err != nil {
		return cluster_info.ClusterInfo{}, err
	}
	ret.Location = loc
	return ret, err
}

func (tr AksTransformer) HubToCloud(hub cluster_info.ClusterInfo, outputScope string) (cluster_info.ClusterInfo, error) {
	var ret = hub
	ret.SourceCluster = &hub
	if hub.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.AZURE
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = outputScope
	loc, err := tr.LocationHubToCloud(ret.Location)
	if err != nil {
		return cluster_info.ClusterInfo{}, errors.Wrap(err, "")
	}
	ret.Location = loc
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
	fn := util.RootPath() + "/locations/azure_locations.csv"
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
