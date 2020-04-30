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

func TranformGCPToHub(clusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	//todo this is duplicate to TransformAzureToHub
	var ret = clusterInfo
	ret.SourceCluster = &clusterInfo
	if clusterInfo.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.HUB
	//	ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in converstion cross-cloud
	loc, err := TransformLocationGcpToHub(ret.Location)
	ret.Location = loc
	return ret, err
}

func TransformHubToGCP(clusterInfo cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	//todo this is duplicate to TransformAzureToHub
	var ret = clusterInfo
	ret.SourceCluster = &clusterInfo
	if clusterInfo.SourceCluster == ret.SourceCluster {
		panic("Copying didn't work as expected")
	}
	ret.Cloud = cluster_info.GCP
	//	ret.Name unchanged
	// ret.NodeCount unchanged
	ret.Scope = "" //Scope not meaningful in conversion cross-cloud
	loc, err := TransformLocationHubToToGcp(ret.Location)
	ret.Location = loc
	return ret, err
}

func TransformLocationGcpToHub(loc string) (string, error) {
	locs, err := getGcpLocations()
	if err != nil {
		return "", err
	}
	if !contains(locs, loc) {
		return "", errors.New(fmt.Sprintf("%s is not a legal location for GCP", loc))
	}
	return loc, nil

}

func TransformLocationHubToToGcp(location string) (string, error) {
	return TransformLocationGcpToHub(location) //locations are taken from GCP, so no conversion; reusing existing code
}
func contains(slice []string, elem string) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}

func getGcpLocations() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PWD", dir)
	ret := make([]string, 20, 20)
	fn := util.RootPath() + "/locations/gcp_locations.csv"
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
		loc := record[0]
		if loc != "" {
			ret = append(ret, loc)
		}
	}
	return ret, nil
}
