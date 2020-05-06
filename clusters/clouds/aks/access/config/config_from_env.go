package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	// these are our *global* config settings, to be shared by all packages.
	// each has corresponding public accessors below.
	// if anything requires a `Set` accessor, that indicates it perhaps
	// shouldn't be set here, because mutable vars shouldn't be global.
	clientID        string
	clientSecret    string
	tenantID        string
	subscriptionID  string
	locationDefault string
	cloudName       = "AzurePublicCloud"
	keepResources   bool
	baseGroupName   string
	userAgent       string
	environment     *azure.Environment
)

// ClientID is the OAuth client ID.
func ClientID() string {
	return clientID
}

// ClientSecret is the OAuth client secret.
func ClientSecret() string {
	return clientSecret
}

// TenantID is the AAD tenant to which this client belongs.
func TenantID() string {
	return tenantID
}

// SubscriptionID is a target subscription for Azure resources.
func SubscriptionID() string {
	return subscriptionID
}

//
// DefaultLocation returns the default location wherein to create new resources.
// Some resource types are not available in all locations so another location might need
// to be chosen.
func DefaultLocation() string {
	return locationDefault
}

// SetBaseGroupName ...
func SetBaseGroupName(name string) {
	baseGroupName = name
}

// BaseGroupName returns a prefix for new groups.
func BaseGroupName() string {
	return baseGroupName
}

// KeepResources specifies whether to keep resources created by samples.
func KeepResources() bool {
	return keepResources
}

// UserAgent specifies a string to append to the agent identifier.
func UserAgent() string {
	if len(userAgent) > 0 {
		return userAgent
	}
	return "sdk-samples"
}

// Environment returns an `azure.Environment{...}` for the current cloud.
func Environment() *azure.Environment {
	if environment != nil {
		return environment
	}
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		// TODO: move to initialization of var
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue\n", cloudName))
	}
	environment = &env
	return environment
}

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
	log.Print("Read Environment")
	return nil
}
