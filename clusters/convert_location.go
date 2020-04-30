package clusters

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func transformLocationAzureToHub(loc string) (string, error) {
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
	fn := "../locations/azure_locations.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		log.Println(err)
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
func transformLocationHubToAzure(location string) (string, error) {
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

func transformLocationGcpToHub(loc string) (string, error) {
	locs, err := getGcpLocations()
	if err != nil {
		return "", err
	}
	if !contains(locs, loc) {
		return "", errors.New(fmt.Sprintf("%s is not a legal location for GCP", loc))
	}
	return loc, nil

}

func transformLocationHubToToGcp(location string) (string, error) {
	return transformLocationGcpToHub(location) //locations are taken from GCP, so no conversion; reusing existing code
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
	fn := "../locations/gcp_locations.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		log.Println("Couldn't open the csv file ", fn, err)
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