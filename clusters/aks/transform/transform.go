package transform

import (
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func TransformAzureToHub(clusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	var ret = clusterInfo
	ret.SourceCluster = &clusterInfo
	if clusterInfo.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.HUB
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in conversion cross-cloud
	loc, err := TransformLocationAzureToHub(ret.Location)
	ret.Location = loc
	return ret, err
}

func TransformHubToAzure(clusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	//todo this is duplicate to TransformAzureToHub
	var ret = clusterInfo
	ret.SourceCluster = &clusterInfo
	if clusterInfo.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.AZURE
	// ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in conversion cross-cloud
	loc, err := TransformLocationHubToAzure(ret.Location)
	ret.Location = loc
	return ret, err
}

//todo split this into Azure and GCP packages
func TransformLocationAzureToHub(loc string) (string, error) {
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
func TransformLocationHubToAzure(location string) (string, error) {
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
	n := make(map[string]string)
	for k, v := range m {
		existing, ok := n[v]
		if !ok || k < existing { //map may not be 1-to-1. If so, take the lexically lowest key as new value
			n[v] = k
		}
	}
	return n
}
