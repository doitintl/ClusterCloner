package config

import (
	"errors"
	"log"
	"os"
	"strconv"
)

// ParseEnvironment loads a sibling `.env` file then looks through all environment
// variables to set global configuration.
func ParseEnvironment() error {
	var err error

	// Use AZURE_BASE_GROUP_NAME and `config.GenerateGroupName()`
	baseGroupName = os.Getenv("AZURE_BASE_GROUP_NAME")

	if baseGroupName == "" {
		err = errors.New("need AZURE_BASE_GROUP_NAME")
		return err
	}
	locationDefault = os.Getenv("AZURE_LOCATION_DEFAULT")
	if locationDefault == "" {
		err = errors.New("need AZURE_LOCATION_DEFAULT")
		return err
	}

	keepResources, err = strconv.ParseBool(os.Getenv("AZURE_SAMPLES_KEEP_RESOURCES"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_SAMPLES_KEEP_RESOURCES, discarding\n")
		keepResources = false
	}

	clientID = os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		err = errors.New("need AZURE_CLIENT_ID")
		return err
	}
	// clientSecret
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		err = errors.New("need AZURE_CLIENT_SECRET")
		return err
	}
	// tenantID (AAD)
	tenantID = os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		err = errors.New("need AZURE_TENANT_ID")
		return err
	}
	// subscriptionID (ARM)
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")

	return nil
}
