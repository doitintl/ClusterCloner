package aks_utils

import (
	"clusterCloner/poc/azure/aks_utils/config"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

// PrintAndLog writes to stdout and to a logger.
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

func LogAndPanic(err error) {
	PrintAndLog(err.Error())
	panic(err)
}

func ReadEnv() error {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	err = addEnv()
	if err != nil {
		return err
	}

	return nil
}

func addEnv() error {
	err := config.ParseEnvironment()
	if err != nil {
		return fmt.Errorf("failed to add top-level env: %+v", err)
	}
	return nil
}
