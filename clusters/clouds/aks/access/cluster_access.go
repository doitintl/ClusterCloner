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

//AksClusterAccess
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

func createGroup(ctx context.Context, groupName string, region string) (resources.Group, error) {
	groupsClient := getGroupsClient()
	log.Println(fmt.Sprintf("Creating resource group '%s' on location: %v", groupName, region))
	return groupsClient.CreateOrUpdate(
		ctx,
		groupName,
		resources.Group{
			Location: to.StringPtr(config.DefaultLocation()),
		})
}

//CreateCluster...
func (ca AksClusterAccess) CreateCluster(createThis cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	grpName := createThis.Scope
	log.Printf("Create Cluster: Group %s, Cluster %s, Location %s", grpName, createThis.Name, createThis.Location)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()

	_, err := createGroup(ctx, grpName, createThis.Location)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "already exists") {
			log.Printf("Group %s already exists", grpName)
		} else {
			return cluster_info.ClusterInfo{}, err
		}
	}
	_, err = createAKSCluster(ctx, createThis.Name, createThis.Location, grpName, aksUsername, aksSSHPublicKeyPath, config.ClientID(), config.ClientSecret(), createThis.K8sVersion, createThis.NodeCount)
	if err != nil {
		log.Println(err)
		return cluster_info.ClusterInfo{}, err
	}
	created := createThis
	created.GeneratedBy = cluster_info.CREATED

	log.Println("Retrieved AKS cluster")
	return created, nil
}

// createAKSCluster creates a new managed Kubernetes cluster
func createAKSCluster(ctx context.Context, resourceName, location, resourceGroupName, username, sshPublicKeyPath, clientID, clientSecret, k8sVersion string, agentPoolCount int32) (c containerservice.ManagedCluster, err error) {
	var sshKeyData string
	if _, err = os.Stat(sshPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
		if err != nil {
			log.Fatalf("Failed to read SSH key data: %v", err)
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
	managedCluster := containerservice.ManagedCluster{
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
			KubernetesVersion:       &k8sVersion,
		},
	}
	future, err := aksClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resourceName,
		managedCluster,
	)
	if err != nil {
		return c, fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	log.Println("About to create Azure Cluster; wait for completion")
	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return c, fmt.Errorf("cannot get the AKS cluster create or update future response: %v", err)
	}

	return future.Result(aksClient)
}

//ListClusters ...
func (ca AksClusterAccess) ListClusters(subscription string, location string) (ci []cluster_info.ClusterInfo, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	var aksClient, err2 = getAKSClient()
	if err2 != nil {
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

		ci := cluster_info.ClusterInfo{
			Scope:       subscription,
			Location:    location,
			Name:        *managedCluster.Name,
			K8sVersion:  *props.KubernetesVersion, //todo could get current version
			NodeCount:   count,
			GeneratedBy: cluster_info.READ,
			Cloud:       cluster_info.AZURE,
		}
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
