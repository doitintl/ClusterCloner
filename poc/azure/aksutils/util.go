package aksutils

import (
	"clustercloner/poc/azure/aksutils/config"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
)

// PrintAndLog writes to stdout and to a logger.
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

// LogAndPanic ...
func LogAndPanic(err error) {
	PrintAndLog(err.Error())
	panic(err)
}

// ReadEnv ...
func ReadEnv() error {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	err = addEnv()
	if err != nil {
		return errors.Wrap(err, "")
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
