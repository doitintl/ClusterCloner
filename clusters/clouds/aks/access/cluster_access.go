package access

import (
	"clustercloner/clusters/clouds/aks/access/config"
	"clustercloner/clusters/clouds/aks/access/iam"
	"clustercloner/clusters/clusterinfo"
	clusterutil "clustercloner/clusters/util"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strconv"

	//"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice" //todo upgrade API

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
	_ = clusterutil.ReadEnv()
}

//AKSClusterAccess ...
type AKSClusterAccess struct {
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

//CreateCluster ...
func (ca AKSClusterAccess) CreateCluster(createThis *clusterinfo.ClusterInfo) (*clusterinfo.ClusterInfo, error) {
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
			return nil, err
		}
	}
	var agentPoolCount int32 = 1
	_, err = createAKSCluster(ctx,
		createThis.Name,
		createThis.Location,
		grpName, aksUsername,
		aksSSHPublicKeyPath,
		config.ClientID(),
		config.ClientSecret(),
		createThis.K8sVersion,
		agentPoolCount,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	created := createThis
	created.GeneratedBy = clusterinfo.CREATED

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
func (ca AKSClusterAccess) ListClusters(subscription string, location string) (ci []*clusterinfo.ClusterInfo, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	var aksClient, err2 = getAKSClient()
	if err2 != nil {
		return ci, errors.New("cannot get AKS client")
	}
	ret := make([]*clusterinfo.ClusterInfo, 0)

	clusterList, _ := aksClient.List(ctx)
	for _, managedCluster := range clusterList.Values() {
		var props = managedCluster.ManagedClusterProperties

		foundCluster := &clusterinfo.ClusterInfo{
			Scope:       subscription,
			Location:    location,
			Name:        *managedCluster.Name,
			K8sVersion:  *props.KubernetesVersion,
			GeneratedBy: clusterinfo.READ,
			Cloud:       clusterinfo.AZURE,
		}

		for _, agentPoolProfile := range *props.AgentPoolProfiles {
			nodePool := clusterinfo.NodePoolInfo{
				Name:        *agentPoolProfile.Name,
				NodeCount:   *agentPoolProfile.Count,
				MachineType: ParseMachineType(fmt.Sprintf("%v", agentPoolProfile.VMSize)),
				DiskSizeGB:  *agentPoolProfile.OsDiskSizeGB,
				K8sVersion:  "",
			}
			foundCluster.AddNodePool(nodePool)
			zero := clusterinfo.MachineType{}
			if nodePool.MachineType == zero {
				panic("cannot read " + agentPoolProfile.VMSize)
			}
		}
		ret = append(ret, foundCluster)

	}
	return ret, nil
}

// ParseMachineType ...
func ParseMachineType(machineType string) clusterinfo.MachineType {
	return MachineTypesNoPromo[machineType] //return zero object if not found
}

// MachineTypes ...
var MachineTypes map[string]clusterinfo.MachineType

// MachineTypesNoPromo ...
var MachineTypesNoPromo map[string]clusterinfo.MachineType

func init() {
	MachineTypes, _ = loadMachineTypes()
	MachineTypesNoPromo = make(map[string]clusterinfo.MachineType)
	for k, v := range MachineTypes {
		if !strings.HasSuffix(k, "Promo") {
			MachineTypesNoPromo[k] = v
		}
	}
}

func loadMachineTypes() (map[string]clusterinfo.MachineType, error) {
	ret := make(map[string]clusterinfo.MachineType)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PWD", dir)

	fn := clusterutil.RootPath() + "/machine-types/aks-vm-sizes.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		wd, _ := os.Getwd()
		log.Println("At ", wd, ":", err)
		return nil, err
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
		cpuFloat, err := strconv.ParseFloat(cpus, 32)
		if err != nil {
			return nil, err
		}
		cpuInt := int32(cpuFloat)

		ramMBString := record[2]
		ramMBFloat, err := strconv.ParseFloat(ramMBString, 32)
		if err != nil {
			return nil, err
		}
		ramGBFloat := ramMBFloat / 1000
		ramGBInt := int32(ramGBFloat)

		ret[name] = clusterinfo.MachineType{Name: name, CPU: cpuInt, RAMGB: ramGBInt}
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
