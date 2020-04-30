package util

import (
	"clusterCloner/clusters/aks/access/config"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadEnv() error {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	if err = config.ParseEnvironment(); err != nil {
		log.Print("Error parsing environment", err)
		return err
	}
	return nil
}

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
