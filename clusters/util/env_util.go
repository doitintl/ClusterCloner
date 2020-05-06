package util

import (
	"clustercloner/clusters/clouds/aks/access/config"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ReadEnv ...
func ReadEnv() error {
	var err error

	rootPath := RootPath()
	envFile := rootPath + "/.env"
	if err := godotenv.Load(envFile); err != nil {
		log.Print("No .env file found at ", envFile)
	}
	if err = config.ParseEnvironment(); err != nil {
		log.Print("Error parsing environment: ", err)
		return errors.Wrap(err, "")

	}
	return nil
}

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
