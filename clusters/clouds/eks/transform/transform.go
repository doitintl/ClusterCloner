package transform

import (
	"clustercloner/clusters"
	clusterutil "clustercloner/clusters/util"
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

// EKSTransformer ...
type EKSTransformer struct{}

// CloudToHub ...
func (tr *EKSTransformer) CloudToHub(in *clusters.ClusterInfo) (ret *clusters.ClusterInfo, err error) {
	panic("")
	return ret, err
}

// HubToCloud ...
func (tr *EKSTransformer) HubToCloud(in *clusters.ClusterInfo, outputScope string) (ret *clusters.ClusterInfo, err error) {
	panic("")
	return ret, err
}

//LocationCloudToHub ...
func (*EKSTransformer) LocationCloudToHub(loc string) (hubValue string, err error) {
	panic("")
	return hubValue, nil
}

//LocationHubToCloud ...
func (EKSTransformer) LocationHubToCloud(location string) (ret string, err error) {
	panic("")
	return ret, nil

}

// Locations ...
var locations map[string]string

// LocationsCloudToHub ...
func LocationsCloudToHub() (map[string]string, error) {
	if locations == nil {
		locations = make(map[string]string)
		fn := clusterutil.RootPath() + "/locations/aws_locations.csv"
		csvfile, err := os.Open(fn)
		if err != nil {
			wd, _ := os.Getwd()
			log.Println("At ", wd, ":", err)
			return nil, err
		}

		r := csv.NewReader(csvfile)
		r.Comma = ';'
		r.Comment = '#'
		first := true
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return nil, errors.Wrap(err, "cannot read line")

			}
			if first {
				first = false
				continue
			}
			if len(record) == 1 {
				log.Println("Short record ", record)
			}
			awsRegion := record[1]
			hubRegion := record[2]
			locations[awsRegion] = hubRegion
		}
	}
	return locations, nil
}
