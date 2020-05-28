package util

import (
	"clustercloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

// LoadLocationMap ...
func LoadLocationMap(file string) (map[string]string, error) {
	ret := make(map[string]string)
	filePath := "/locations/" + file
	fn := util.RootPath() + filePath
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
			return nil, errors.Wrap(err, "cannot read "+file)
		}
		if first {
			first = false
			continue
		}
		if len(record) == 1 {
			log.Println("Short record", record)
		}
		cloudRegion := record[3]
		hubRegion := record[5]
		supportsAks := record[4]
		if supportsAks != "true" {
			return nil, errors.New(fmt.Sprintf("Azure region %s file not supported in %s", cloudRegion, file))
		}
		ret[cloudRegion] = hubRegion
	}
	return ret, nil
}
