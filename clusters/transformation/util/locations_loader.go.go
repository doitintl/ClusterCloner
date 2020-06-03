package util

import (
	"clustercloner/clusters/util"
	"encoding/csv"
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
		return nil, errors.Wrap(err, "error opening "+fn)
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
			return nil, errors.Wrap(err, "cannot read "+file)
		}
		if first {
			first = false
			continue
		}
		if len(record) == 1 {
			log.Println("Short record", record)
		}
		if len(record) != 3 {
			return nil, errors.Errorf("wrong length record, length %d (%v)", len(record), record)
		}
		cloudRegion := record[1]
		hubRegion := record[2]
		ret[cloudRegion] = hubRegion
	}
	return ret, nil
}
