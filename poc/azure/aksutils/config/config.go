package config

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	// these are our *global* config settings, to be shared by all packages.
	// each has corresponding public accessors below.
	// if anything requires a `Set` accessor, that indicates it perhaps
	// shouldn't be set here, because mutable vars shouldn't be global.
	clientID               string
	clientSecret           string
	tenantID               string
	subscriptionID         string
	locationDefault        string
	authorizationServerURL string
	cloudName              = "AzurePublicCloud"
	keepResources          bool
	baseGroupName          string
	userAgent              string
	environment            *azure.Environment
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

// DefaultLocation  returns the default location wherein to create new resources.
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
