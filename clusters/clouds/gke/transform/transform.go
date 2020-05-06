package transform

import (
	"clustercloner/clusters/clusterinfo"
	transformutil "clustercloner/clusters/transformation/util"
	"clustercloner/clusters/util"
	"encoding/csv"
	"github.com/pkg/errors"

	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

var randNumGen *rand.Rand

func init() {
	s := rand.NewSource(time.Now().Unix())
	randNumGen = rand.New(s) // initialize local pseudorandom generator
}

// GkeTransformer ...
type GkeTransformer struct {
}

// CloudToHub ...
func (tr *GkeTransformer) CloudToHub(in clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error) {
	loc, err := tr.LocationCloudToHub(in.Location)
	if err != nil {
		return clusterinfo.ClusterInfo{}, errors.Wrap(err, "error in converting locations")
	}
	k8sVersion, err := transformutil.MajorMinorPatchVersion(in.K8sVersion)
	if err != nil {
		return clusterinfo.ClusterInfo{}, errors.Wrap(err, "error in K8s Version "+in.K8sVersion)
	}

	ret := transformutil.TransformSpoke(in, "", clusterinfo.HUB, loc, k8sVersion)
	return ret, err
}

// HubToCloud ...
func (tr *GkeTransformer) HubToCloud(in clusterinfo.ClusterInfo, outputScope string) (clusterinfo.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return clusterinfo.ClusterInfo{}, errors.Wrap(err, "error in converting location")
	}
	ret := transformutil.TransformSpoke(in, outputScope, clusterinfo.GCP, loc, in.K8sVersion)

	return ret, err
}

// LocationCloudToHub ...
func (tr *GkeTransformer) LocationCloudToHub(zone string) (string, error) {
	locs, err := getGcpLocations()
	if err != nil {
		return "", err
	}
	hyphenCount, secondHyphenIdx := Hyphens(zone)
	if hyphenCount != 1 && hyphenCount != 2 {
		msg := fmt.Sprintf("%s is not a legal zone/region format for GCP", zone)
		log.Println(msg)
		return "", errors.New(msg)
	}
	runes := []rune(zone)
	endRegion := len(runes)
	if secondHyphenIdx > 1 {
		endRegion = secondHyphenIdx
	}
	region := string(runes[0:endRegion])
	if !contains(locs, region) {
		msg := fmt.Sprintf("Zone %s is not in a legal region for GCP", zone)
		log.Println(msg)
		return "", errors.New(msg)
	}
	return region, nil

}

// Hyphens ...
func Hyphens(zone string) (hyphenCount int, secondHyphenIdx int) {
	secondHyphenIdx = -1
	for i, ch := range zone {
		if ch == '-' {
			hyphenCount++
			if hyphenCount == 2 {
				secondHyphenIdx = i
			}
		}
	}
	return hyphenCount, secondHyphenIdx
}

// LocationHubToCloud ...
func (GkeTransformer) LocationHubToCloud(location string) (string, error) {
	hyphenCount, _ := Hyphens(location)
	var zone string
	if hyphenCount == 1 {
		zones := []string{"a", "b"}
		//Even when converting GCP to GCP, use a random zone, because we decided to convert GCP to GCP through the Hub format.
		var randIdx = randNumGen.Intn(len(zones))
		randZone := zones[randIdx]
		zone = location + "-" + randZone
	} else if hyphenCount == 2 {
		zone = location
	} else {
		panic(location)
	}
	return zone, nil

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
