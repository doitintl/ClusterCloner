package aks

import (
	"clusterCloner/clusters/aks/utils/config"
	"clusterCloner/clusters/aks/utils/iam"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/go-autorest/autorest/to"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func getAKSClient() (containerservice.ManagedClustersClient, error) {
	aksClient := containerservice.NewManagedClustersClient(config.SubscriptionID())
	auth, _ := iam.GetResourceManagementAuthorizer()
	aksClient.Authorizer = auth
	aksClient.AddToUserAgent(config.UserAgent())
	aksClient.PollingDuration = time.Hour * 1
	return aksClient, nil
}

// createAKSCluster creates a new managed Kubernetes cluster
func createAKSCluster(ctx context.Context, resourceName, location, resourceGroupName, username, sshPublicKeyPath, clientID, clientSecret string, agentPoolCount int32) (c containerservice.ManagedCluster, err error) {
	var sshKeyData string
	if _, err = os.Stat(sshPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
		if err != nil {
			log.Fatalf("failed to read SSH key data: %v", err)
		}
		sshKeyData = string(sshBytes)
	} else {
		log.Printf("Cannot load: %s", sshPublicKeyPath)
		sshKeyData = "fakepubkey"
	}

	aksClient, err := getAKSClient()
	if err != nil {
		return c, fmt.Errorf("cannot get AKS client: %v", err)
	}

	future, err := aksClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resourceName,
		containerservice.ManagedCluster{
			Name:     &resourceName,
			Location: &location,
			ManagedClusterProperties: &containerservice.ManagedClusterProperties{
				DNSPrefix: &resourceName,
				LinuxProfile: &containerservice.LinuxProfile{
					AdminUsername: to.StringPtr(username),
					SSH: &containerservice.SSHConfiguration{
						PublicKeys: &[]containerservice.SSHPublicKey{
							{
								KeyData: to.StringPtr(sshKeyData),
							},
						},
					},
				},
				AgentPoolProfiles: &[]containerservice.AgentPoolProfile{
					{
						Count:  to.Int32Ptr(agentPoolCount),
						Name:   to.StringPtr("agentpool1"),
						VMSize: containerservice.StandardD2sV3,
					},
				},
				ServicePrincipalProfile: &containerservice.ServicePrincipalProfile{
					ClientID: to.StringPtr(clientID),
					Secret:   to.StringPtr(clientSecret),
				},
			},
		},
	)
	if err != nil {
		return c, fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return c, fmt.Errorf("cannot get the AKS cluster create or update future response: %v", err)
	}

	return future.Result(aksClient)
}
