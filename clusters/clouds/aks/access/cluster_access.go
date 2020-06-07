package access

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clusteraccess/util"
	"clustercloner/clusters/machinetypes"
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
	defer clusterutil.TrackTime("Delete AKS", time.Now())

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	aksClient, err := getManagedClustersClient()
	if err != nil {
		return errors.New("cannot get AKS client")
	}
	future, err := aksClient.Delete(ctx, deleteThis.Scope, deleteThis.Name)
	if err != nil {
		return fmt.Errorf("cannot delete cluster: %v", err)
	}

	log.Printf("About to delete AKS Cluster %s; waiting for completion", deleteThis.Name)
	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the AKS deletion  future response: %v", err)
	}
	response, err := future.Result(aksClient)
	if err != nil {
		return errors.Wrap(err, "error waiting for result")
	}
	status := response.StatusCode
	if status != 200 {
		return errors.New("could not delete cluster, state was " + response.Status)
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
	group := resources.Group{Location: &region}
	return groupsClient.CreateOrUpdate(ctx, groupName, group)
}

//Describe ...
func (ca AKSClusterAccess) Describe(searchTemplate *clusters.ClusterInfo) (described *clusters.ClusterInfo, err error) {
	defer clusterutil.TrackTime("Describe AKS", time.Now())
	groupName := searchTemplate.Scope
	name := searchTemplate.Name
	log.Println("Describe AKS cluster", groupName, ": ", name)
	if searchTemplate.GeneratedBy == "" {
		searchTemplate.GeneratedBy = clusters.SearchTemplate
	}
	if searchTemplate.GeneratedBy != clusters.SearchTemplate {
		log.Printf("Wrong GeneratedBy: %s\n", searchTemplate.GeneratedBy)
	}

	cluster, err := getCluster(groupName, name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.Wrap(err, "cluster "+name+"not found")
		}
		return nil, errors.Wrap(err, "cannot get cluster "+name)
	}
	clusterInfo, err := clusterObjectToClusterInfo(cluster, searchTemplate.Scope, clusters.Read)
	if err != nil {
		return nil, errors.New("cannot convert cluster object for " + name)
	}

	clusterInfo.SourceCluster = searchTemplate
	return clusterInfo, nil
}

// Create ...
func (ca AKSClusterAccess) Create(createThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error) {
	defer clusterutil.TrackTime("Create AKS", time.Now())

	groupName := createThis.Scope
	log.Printf("Create AKS Cluster: Group %s, Name %s, Location %s", groupName, createThis.Name, createThis.Location)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	_, err = createGroup(ctx, groupName, createThis.Location)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "already exists, proceeding") ||
			(strings.Contains(errS, "Invalid resource group location") &&
				strings.Contains(errS, "The Resource group already exists in location")) {
			log.Printf("Group %s already exists; no need to create", groupName)
		} else {
			return nil, errors.Wrap(err, "error creating group")
		}
	} else {
		log.Println("Created resource group " + groupName)
	}

	createdCluster, err := createAKSCluster(ctx, createThis, aksUsername, aksSSHPublicKeyPath, ClientID(), ClientSecret())
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot create cluster")
	}
	createdClusterInfo, err := clusterObjectToClusterInfo(createdCluster, createThis.Scope, clusters.Created)
	if err != nil {
		return nil, errors.New("cannot convert cluster object")
	}

	createdClusterInfo.SourceCluster = createThis
	return createdClusterInfo, nil
}

// createAKSCluster creates a new managed Kubernetes cluster
func createAKSCluster(ctx context.Context,
	createThis *clusters.ClusterInfo,
	username, sshPublicKeyPath, clientID, clientSecret string) (created containerservice.ManagedCluster, err error) {
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
	for _, npi := range createThis.NodePools {

		var scaleSetPriority containerservice.ScaleSetPriority
		if npi.Preemptible {
			scaleSetPriority = containerservice.Spot
		} else {
			scaleSetPriority = containerservice.Regular
		}
		agPoolName := strings.ReplaceAll(npi.Name, "-", "")
		agPoolProfile := containerservice.ManagedClusterAgentPoolProfile{
			Count:            to.Int32Ptr(int32(npi.NodeCount)),
			Name:             to.StringPtr(agPoolName),
			VMSize:           containerservice.VMSizeTypes(npi.MachineType.Name),
			OsDiskSizeGB:     to.Int32Ptr(int32(npi.DiskSizeGB)),
			ScaleSetPriority: scaleSetPriority,
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
		createThis.Scope,
		createThis.Name,
		clusterToCreate,
	)
	if err != nil {
		return containerservice.ManagedCluster{}, fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	log.Printf("About to create Azure Cluster %s; waiting for completion", createThis.Name)
	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return containerservice.ManagedCluster{}, fmt.Errorf("cannot WaitForCompletion on the response from CreateOrUpdate: %v", err)

	}
	created, err = future.Result(aksClient)
	if err != nil {
		return containerservice.ManagedCluster{}, errors.Wrap(err, "error waiting for result")
	}
	clusterProperties := created.ManagedClusterProperties
	state := *clusterProperties.ProvisioningState
	if state != "Succeeded" {
		return containerservice.ManagedCluster{}, errors.New("could not created cluster, state was " + state)
	}
	return created, nil
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
		foundCluster, err := clusterObjectToClusterInfo(managedCluster, subscription, clusters.Read)
		if err != nil {
			return nil, errors.New("cannot convert cluster object")
		}

		ret = append(ret, foundCluster)

	}
	util.PrintFilteringResults(clusters.Azure, labelFilter, matchedNames, unmatchedNames)

	return ret, nil
}

func clusterObjectToClusterInfo(managedCluster containerservice.ManagedCluster, subscription string, generatedBy string) (*clusters.ClusterInfo, error) {
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
	agentPoolProfilesPtr := props.AgentPoolProfiles
	var agentPoolProfiles []containerservice.ManagedClusterAgentPoolProfile = nil
	if agentPoolProfilesPtr != nil {
		agentPoolProfiles = *agentPoolProfilesPtr
	}
	for _, agentPoolProfile := range agentPoolProfiles {
		var scaleSetPriority = agentPoolProfile.ScaleSetPriority
		var spot = scaleSetPriority == containerservice.Spot
		machType, err := aksMachineTypes.Get(fmt.Sprintf("%v", agentPoolProfile.VMSize))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("cannot get machine type %v", agentPoolProfile.VMSize))
		}
		if machType.Name == "" {
			return nil, errors.Errorf("cannot find machine type %v", agentPoolProfile.VMSize)
		}
		npi := clusters.NodePoolInfo{
			Name:        *agentPoolProfile.Name,
			NodeCount:   int(*agentPoolProfile.Count),
			MachineType: machType,
			DiskSizeGB:  int(*agentPoolProfile.OsDiskSizeGB),
			K8sVersion:  nodePoolK8sVersion,
			Preemptible: spot,
		}
		foundCluster.AddNodePool(npi)
		zero := machinetypes.MachineType{}
		if npi.MachineType == zero {
			panic("cannot read " + agentPoolProfile.VMSize)
		}
	}
	return foundCluster, nil
}

// supportedVersions ...
var supportedVersions []string

// GetSupportedK8sVersions ...
func (ca AKSClusterAccess) GetSupportedK8sVersions(scope, location string) ([]string, error) {

	if supportedVersions == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
		defer cancel()
		supportedVersions = make([]string, 0)

		listOrch, err := getContainerServicesClient().ListOrchestrators(ctx, location, "")
		if err != nil {
			return nil, errors.Wrap(err, "cannot ListOrchestrators for "+location)
		}
		for _, orch := range *listOrch.Orchestrators {
			t := *orch.OrchestratorType
			if t == "Kubernetes" {
				supportedVersions = append(supportedVersions, *orch.OrchestratorVersion)
			}

		}
	}
	return supportedVersions, nil
}

// aksMachineTypes ...
var aksMachineTypes *machinetypes.MachineTypes

// GetMachineTypes ...
func GetMachineTypes() *machinetypes.MachineTypes {
	return aksMachineTypes
}

func init() {
	var err error
	aksMachineTypes, err = loadMachineTypes()
	if err != nil {
		panic(fmt.Sprintf("cannot load AKS machine types %v", err))
	}
	if aksMachineTypes == nil || aksMachineTypes.Length() == 0 {
		panic(fmt.Sprintf("cannot load AKS machine types %v", err))
	}
}

func loadMachineTypes() (*machinetypes.MachineTypes, error) {
	ret := machinetypes.NewMachineTypeMap()

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
			return nil, errors.Wrap(err, "cannot parse "+ramMBString)
		}
		ret.Set(name, machinetypes.MachineType{Name: name, CPU: int(cpuInteger), RAMMB: int(ramMBInt)})
	}
	return &ret, nil
}

//
func getManagedClustersClient() (mcc containerservice.ManagedClustersClient, err error) {
	client := containerservice.NewManagedClustersClient(SubscriptionID())
	auth, err := GetResourceManagementAuthorizer()
	if err != nil {
		return mcc, errors.Wrap(err, "cannot GetResourceManagementAuthorizer")
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
