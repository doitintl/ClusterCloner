package util

import (
	"clusterCloner/clusters/aks/config"
	"github.com/joho/godotenv"
	"log"
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
