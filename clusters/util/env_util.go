package util

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// RootPath ...
func RootPath() string {
	wd, _ := os.Getwd()
	for wd != "" {
		files, err := ioutil.ReadDir(wd)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if "locations" == f.Name() {
				return wd
			}
		}
		if wd == "/" {
			wd = ""
		} else {
			wd = filepath.Dir(wd)
		}
	}
	return wd
}

// NoopWriter ...
type NoopWriter struct{}

// Write ...
func (m NoopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
