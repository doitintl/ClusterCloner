package access

import (
	"clusterCloner/clusters/clouds/aks/access/config"
	"clusterCloner/clusters/clouds/aks/access/iam"
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"context"
	"errors"
	"fmt"
	//	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/to"

	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	aksUsername         = "azureuser"
	aksSSHPublicKeyPath = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
)

func init() {
	_ = util.ReadEnv()
}

type AksClusterAccess struct {
}

func getGroupsClient() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(config.SubscriptionID())
	a, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		log.Fatalf("failed to initialize authorizer: %v\n", err)
	}
	groupsClient.Authorizer = a
	_ = groupsClient.AddToUserAgent(config.UserAgent())
	return groupsClient
}

// createGroup creates a new resource group named by env var
func createGroup(ctx context.Context, groupName string) (resources.Group, error) {
	groupsClient := getGroupsClient()
	log.Println(fmt.Sprintf("creating resource group '%s' on location: %v", groupName, config.DefaultLocation()))
	return groupsClient.CreateOrUpdate(
		ctx,
		groupName,
		resources.Group{
			Location: to.StringPtr(config.DefaultLocation()),
		})
}
func (ca AksClusterAccess) CreateCluster(info cluster_info.ClusterInfo) error {
	grpName := config.BaseGroupName()
	log.Printf("Group %s, Cluster %s, Location %s", grpName, info.Name, info.Location)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()

	_, err := createGroup(ctx, grpName)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "already exists") {
			log.Printf("Group %s already exists", grpName)
			return err
		} else {
			log.Print(err)
			return err
		}
	}
	_, err = createAKSCluster(ctx, info.Name, info.Location, grpName, aksUsername, aksSSHPublicKeyPath, config.ClientID(), config.ClientSecret(), info.NodeCount)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Retrieved RAKS cluster")
	return nil
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

	agentPoolProfiles := &[]containerservice.AgentPoolProfile{
		{
			Count:  to.Int32Ptr(agentPoolCount),
			Name:   to.StringPtr("agentpool1"),
			VMSize: containerservice.StandardD2sV3,
		},
	}
	servicePrincipalProfile := &containerservice.ServicePrincipalProfile{
		ClientID: to.StringPtr(clientID),
		Secret:   to.StringPtr(clientSecret),
	}
	sshConfiguration := &containerservice.SSHConfiguration{
		PublicKeys: &[]containerservice.SSHPublicKey{
			{
				KeyData: to.StringPtr(sshKeyData),
			},
		},
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
					SSH:           sshConfiguration,
				},
				AgentPoolProfiles:       agentPoolProfiles,
				ServicePrincipalProfile: servicePrincipalProfile,
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

func (ca AksClusterAccess) ListClusters(subscription string, location string) (ci []cluster_info.ClusterInfo, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	var aksClient, err_ = getAKSClient()
	if err_ != nil {
		return ci, errors.New("cannot get AKS client")
	}
	ret := make([]cluster_info.ClusterInfo, 0)

	clusterList, _ := aksClient.List(ctx)
	for _, managedCluster := range clusterList.Values() {
		var props = managedCluster.ManagedClusterProperties

		var count int32 = 0
		for _, app := range *props.AgentPoolProfiles {
			count += *app.Count
		}
		ci := cluster_info.ClusterInfo{Scope: subscription, Location: location, Name: *managedCluster.Name, NodeCount: count, GeneratedBy: cluster_info.READ}
		ret = append(ret, ci)

	}
	return ret, nil
}

func getAKSClient() (mcc containerservice.ManagedClustersClient, err error) {
	aksClient := containerservice.NewManagedClustersClient(config.SubscriptionID())
	auth, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return mcc, err
	}
	aksClient.Authorizer = auth
	_ = aksClient.AddToUserAgent(config.UserAgent())
	aksClient.PollingDuration = time.Hour * 1
	return aksClient, nil
}
