package access

import (
	"clustercloner/clusters"
	clusterutil "clustercloner/clusters/util"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/pkg/errors"
	"io"
	"strconv"

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
	err := ReadEnv()
	if err != nil {
		panic("Cannot read environment, aborting")
	}
}

//AKSClusterAccess ...
type AKSClusterAccess struct {
}

// Delete ...
func (ca AKSClusterAccess) Delete(deleteThis *clusters.ClusterInfo) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	aksClient, err := getManagedClustersClient()
	if err != nil {
		return errors.New("cannot get AKS client")
	}
	future, err := aksClient.Delete(ctx, deleteThis.Scope, deleteThis.Name)
	if err != nil {
		return fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	log.Printf("About to delete Azure Cluster %s; wait for completion", deleteThis.Name)
	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the AKS create or update future response: %v", err)
	}
	response, err := future.Result(aksClient)
	if err != nil {
		return errors.Wrap(err, "error waiting for result")
	}
	status := response.StatusCode
	if status != 200 {
		return errors.New("could not created cluster, staate was " + response.Status)
	}
	return nil
}

// GetAKS ...
func getCluster(resourceGroupName, resourceName string) (cluster containerservice.ManagedCluster, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	aksClient, err := getManagedClustersClient()
	if err != nil {
		return cluster, fmt.Errorf("cannot get AKS client: %v", err)
	}

	cluster, err = aksClient.Get(ctx, resourceGroupName, resourceName)
	if err != nil {
		return cluster, fmt.Errorf("cannot get AKS managed cluster %v from resource group %v: %v", resourceName, resourceGroupName, err)
	}

	return cluster, nil
}
func createGroup(ctx context.Context, groupName string, region string) (resources.Group, error) {
	groupsClient := getGroupsClient()
	log.Println(fmt.Sprintf("Creating resource group '%s' on location: %v", groupName, region))
	group := resources.Group{Location: to.StringPtr(DefaultLocation())}
	return groupsClient.CreateOrUpdate(ctx, groupName, group)
}

//Describe ...
func (ca AKSClusterAccess) Describe(searchTemplate *clusters.ClusterInfo) (described *clusters.ClusterInfo, err error) {
	if searchTemplate.GeneratedBy == "" {
		searchTemplate.GeneratedBy = clusters.SearchTemplate
	}
	if searchTemplate.GeneratedBy != clusters.SearchTemplate {
		panic(fmt.Sprintf("Wrong GeneratedBy: %s" + searchTemplate.GeneratedBy))
	}
	groupName := searchTemplate.Scope

	cluster, err := getCluster(groupName, searchTemplate.Name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster")
	}
	clusterInfo := clusterObjectToClusterInfo(cluster, searchTemplate.Scope, clusters.Read)
	clusterInfo.SourceCluster = searchTemplate
	return clusterInfo, nil
}

// Create ...
func (ca AKSClusterAccess) Create(createThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error) {

	groupName := createThis.Scope
	log.Printf("Create Cluster: Group %s, Cluster %s, Location %s", groupName, createThis.Name, createThis.Location)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	_, err = createGroup(ctx, groupName, createThis.Location)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "already exists, proceeding") ||
			(strings.Contains(errS, "Invalid resource group location") &&
				strings.Contains(errS, "The Resource group already exists in location")) {
			log.Printf("Group %s already exists: %v", groupName, err)
		} else {
			return nil, err
		}
	}

	createdCluster, err := createAKSCluster(ctx, createThis, groupName, aksUsername, aksSSHPublicKeyPath, ClientID(), ClientSecret())
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot create cluster")
	}
	createdClusterInfo := clusterObjectToClusterInfo(createdCluster, createThis.Scope, clusters.Created)
	createdClusterInfo.SourceCluster = createThis
	return createdClusterInfo, nil
}

// createAKSCluster creates a new managed Kubernetes cluster
func createAKSCluster(ctx context.Context, createThis *clusters.ClusterInfo,
	resourceGroupName, username, sshPublicKeyPath, clientID, clientSecret string) (clusterCreated containerservice.ManagedCluster, err error) {
	var sshKeyData string
	if _, err := os.Stat(sshPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
		if err != nil {
			log.Fatalf("Failed to read SSH key data: %v", err)
		}
		sshKeyData = string(sshBytes)
	} else {
		panic(fmt.Sprintf("cannot load: %s", sshPublicKeyPath))
	}
	aksClient, err := getManagedClustersClient()
	if err != nil {
		return containerservice.ManagedCluster{}, fmt.Errorf("cannot get AKS client: %v", err)
	}

	agPoolProfiles := make([]containerservice.ManagedClusterAgentPoolProfile, 0)
	for _, nodePool := range createThis.NodePools {
		agPoolName := strings.ReplaceAll(nodePool.Name, "-", "")
		agPoolProfile := containerservice.ManagedClusterAgentPoolProfile{
			Count:        to.Int32Ptr(nodePool.NodeCount),
			Name:         to.StringPtr(agPoolName),
			VMSize:       containerservice.VMSizeTypes(nodePool.MachineType.Name),
			OsDiskSizeGB: to.Int32Ptr(nodePool.DiskSizeGB),
			//TODO use the nodePool.K8sVersion. Does Az support that?
		}
		agPoolProfiles = append(agPoolProfiles, agPoolProfile)
	}

	servicePrincipalProfile := &containerservice.ManagedClusterServicePrincipalProfile{
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
	clusterToCreate := containerservice.ManagedCluster{
		Name:     &createThis.Name,
		Location: &createThis.Location,
		Tags:     clusterutil.StrMapToStrPtrMap(createThis.Labels),
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			DNSPrefix: &createThis.Name,
			LinuxProfile: &containerservice.LinuxProfile{
				AdminUsername: to.StringPtr(username),
				SSH:           sshConfiguration,
			},
			AgentPoolProfiles:       &agPoolProfiles,
			ServicePrincipalProfile: servicePrincipalProfile,
			KubernetesVersion:       &createThis.K8sVersion,
		},
	}
	future, err := aksClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		createThis.Name,
		clusterToCreate,
	)
	if err != nil {
		return containerservice.ManagedCluster{}, fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	log.Printf("About to create Azure Cluster %s; wait for completion", createThis.Name)
	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return containerservice.ManagedCluster{}, fmt.Errorf("cannot get the AKS create or update future response: %v", err)
	}
	clusterCreated, err = future.Result(aksClient)
	if err != nil {
		return containerservice.ManagedCluster{}, errors.Wrap(err, "error waiting for result")
	}
	clusterProperties := clusterCreated.ManagedClusterProperties
	state := *clusterProperties.ProvisioningState
	if state != "Succeeded" {
		return containerservice.ManagedCluster{}, errors.New("could not created cluster, staate was " + state)
	}
	return clusterCreated, err
}

// List ...
func (ca AKSClusterAccess) List(subscription, location string, labelFilter map[string]string) (listedClusters []*clusters.ClusterInfo, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	aksClient, err := getManagedClustersClient()
	if err != nil {
		return nil, errors.New("cannot get AKS client")
	}

	ret := make([]*clusters.ClusterInfo, 0)

	clusterList, err := aksClient.List(ctx)
	// TODO check if we need paging
	if err != nil {
		return nil, errors.Wrap(err, "cannot list")
	}

	unmatchedNames := make([]string, 0)
	matchedNames := make([]string, 0)
	for _, managedCluster := range clusterList.Values() {
		tags := managedCluster.Tags
		tagsAsStrMap := clusterutil.StrPtrMapToStrMap(tags)
		match := clusterutil.LabelMatch(labelFilter, tagsAsStrMap)
		name := *managedCluster.Name
		if !match {
			log.Printf("Skipping cluster %s because labels do not match", name)
			unmatchedNames = append(unmatchedNames, name)
			continue
		}
		matchedNames = append(matchedNames, name)
		foundCluster := clusterObjectToClusterInfo(managedCluster, subscription, clusters.Read)
		ret = append(ret, foundCluster)

	}
	log.Printf("In listing clusters, these matched the label filter %v; and these did not %v\n", matchedNames, unmatchedNames)

	return ret, nil
}

func clusterObjectToClusterInfo(managedCluster containerservice.ManagedCluster, subscription string, generatedBy string) *clusters.ClusterInfo {
	var props = managedCluster.ManagedClusterProperties
	foundCluster := &clusters.ClusterInfo{
		Scope:       subscription,
		Location:    *managedCluster.Location,
		Name:        *managedCluster.Name,
		K8sVersion:  *props.KubernetesVersion,
		GeneratedBy: generatedBy,
		Cloud:       clusters.Azure,
		Labels:      clusterutil.StrPtrMapToStrMap(managedCluster.Tags),
	}
	//AgentPoolProfile is not showing AgentPool K8s Version, so copying from the Cluster
	var nodePoolK8sVersion = foundCluster.K8sVersion
	for _, agentPoolProfile := range *props.AgentPoolProfiles {
		nodePool := clusters.NodePoolInfo{
			Name:        *agentPoolProfile.Name,
			NodeCount:   *agentPoolProfile.Count,
			MachineType: MachineTypeByName(fmt.Sprintf("%v", agentPoolProfile.VMSize)),
			DiskSizeGB:  *agentPoolProfile.OsDiskSizeGB,
			K8sVersion:  nodePoolK8sVersion,
		}
		foundCluster.AddNodePool(nodePool)
		zero := clusters.MachineType{}
		if nodePool.MachineType == zero {
			panic("cannot read " + agentPoolProfile.VMSize)
		}
	}
	return foundCluster
}

// supportedVersions ...
var supportedVersions []string

// GetSupportedK8sVersions ...
func (ca AKSClusterAccess) GetSupportedK8sVersions(scope, location string) []string {

	if supportedVersions == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
		defer cancel()
		supportedVersions = make([]string, 0)

		listOrch, err := getContainerServicesClient().ListOrchestrators(ctx, location, "")
		if err != nil {
			log.Println(err)
		} else {
			for _, orch := range *listOrch.Orchestrators {
				t := *orch.OrchestratorType
				//				log.Println(*orch.OrchestratorType, *orch.OrchestratorVersion)
				if t == "Kubernetes" {
					supportedVersions = append(supportedVersions, *orch.OrchestratorVersion)
				}
			}
		}
	}
	return supportedVersions
}

// MachineTypeByName ...
func MachineTypeByName(machineType string) clusters.MachineType {
	return MachineTypes[machineType] //return zero object if not found
}

// MachineTypes ...
var MachineTypes map[string]clusters.MachineType

func init() {
	var err error
	MachineTypes, err = loadMachineTypes()
	if len(MachineTypes) == 0 || err != nil {
		panic(fmt.Sprintf("cannot load machine types %v", err))
	}
}

func loadMachineTypes() (map[string]clusters.MachineType, error) {
	ret := make(map[string]clusters.MachineType)

	fn := clusterutil.RootPath() + "/machine-types/aks-vm-sizes.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		wd, _ := os.Getwd()
		return nil, errors.Wrap(err, fmt.Sprintf("At %s: %v", wd, err))
	}

	r := csv.NewReader(csvfile)
	r.Comma = ','
	first := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return nil, errors.Wrap(err, "cannot read csv")
		}
		if first {
			first = false
			continue
		}
		if len(record) == 1 {
			log.Println("Short record ", record)
		}
		name := record[0]

		cpus := record[1]
		cpuInteger, err := strconv.ParseInt(cpus, 10, 32)
		if err != nil || cpuInteger == 0 {
			return nil, errors.Wrap(err, "cannot parse cpus "+cpus)
		}

		ramMBString := record[2]
		ramMBInt, err := strconv.ParseInt(ramMBString, 10, 32)
		if err != nil {
			return nil, err
		}
		ret[name] = clusters.MachineType{Name: name, CPU: int32(cpuInteger), RAMMB: int32(ramMBInt)}
	}
	return ret, nil
}

//
func getManagedClustersClient() (mcc containerservice.ManagedClustersClient, err error) {
	client := containerservice.NewManagedClustersClient(SubscriptionID())
	auth, err := GetResourceManagementAuthorizer()
	if err != nil {
		return mcc, err
	}
	client.Authorizer = auth
	_ = client.AddToUserAgent(UserAgent())
	return client, nil
}

func getGroupsClient() resources.GroupsClient {
	client := resources.NewGroupsClient(SubscriptionID())
	auth, err := GetResourceManagementAuthorizer()
	if err != nil {
		log.Fatalf("failed to initialize authorizer: %v\n", err)
	}
	client.Authorizer = auth
	_ = client.AddToUserAgent(UserAgent())
	return client
}
func getContainerServicesClient() containerservice.ContainerServicesClient {

	client := containerservice.NewContainerServicesClient(SubscriptionID())
	auth, err := GetResourceManagementAuthorizer()
	if err != nil {
		log.Fatalf("failed to initialize authorizer: %v\n", err)
	}
	client.Authorizer = auth
	_ = client.AddToUserAgent(UserAgent())
	return client
}
