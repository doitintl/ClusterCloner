package transform

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/gke/access"
	transformutil "clustercloner/clusters/transformation/util"
	baseutil "clustercloner/clusters/util"
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

// GKETransformer ...
type GKETransformer struct {
}

// GKEToGKETransformer ...
type GKEToGKETransformer struct {
	transformutil.IdentityTransformer
}

// CloudToHub ...
func (tr *GKETransformer) CloudToHub(in *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationCloudToHub(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting locations")
	}
	clusterK8sVersion, err := baseutil.MajorMinorPatchVersion(in.K8sVersion)
	if err != nil {
		return nil, errors.Wrap(err, "error in K8s K8sVersion "+in.K8sVersion)
	}

	ret, err := transformutil.TransformSpoke(in, "", clusters.Hub, loc, clusterK8sVersion, nil, false)

	return ret, err
}

// HubToCloud ...
func (tr *GKETransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting location")
	}
	ret, err := transformutil.TransformSpoke(in, outputScope, clusters.GCP, loc, in.K8sVersion, access.MachineTypes, true)

	return ret, err
}

// LocationCloudToHub ...
func (tr *GKETransformer) LocationCloudToHub(zone string) (string, error) {
	locs, err := getGcpLocations()
	if err != nil {
		return "", err
	}
	hyphenCount, secondHyphenIdx := hyphensForGCPLocation(zone)
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
	if !baseutil.ContainsStr(locs, region) {
		msg := fmt.Sprintf("Zone %s is not in a legal region for GCP", zone)
		log.Println(msg)
		return "", errors.New(msg)
	}
	return region, nil

}

// HubToCloud ...
func (tr *GKEToGKETransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (*clusters.ClusterInfo, error) {
	loc, err := tr.LocationHubToCloud(in.Location)
	if err != nil {
		return nil, errors.Wrap(err, "error in converting location")
	}
	ret, err := transformutil.TransformSpoke(in, outputScope, clusters.GCP, loc, in.K8sVersion, access.MachineTypes, true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not TransformSpoke")
	}
	return ret, err
}

// hyphensForGCPLocation ...
func hyphensForGCPLocation(zone string) (hyphenCount int, secondHyphenIdx int) {
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
func (GKETransformer) LocationHubToCloud(location string) (string, error) {
	hyphenCount, _ := hyphensForGCPLocation(location)
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

func getGcpLocations() ([]string, error) {

	ret := make([]string, 20)
	fn := baseutil.RootPath() + "/locations/gcp_locations.csv"
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
