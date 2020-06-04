package access

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/awssdk"
	"clustercloner/clusters/clouds/eks/eksctl"
	accessutil "clustercloner/clusters/clusteraccess/util"
	"clustercloner/clusters/machinetypes"
	clusterutil "clustercloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	key := "AWS_SHARED_CREDENTIALS_FILE"
	cred := os.Getenv(key)
	if cred == "" {
		log.Println("No " + key + " env variable, so using awscredentials in application root as default")
		cred = "awscredentials"
	}
	rootPath := clusterutil.RootPath() + "/" + cred
	err := os.Setenv(key, rootPath)
	if err != nil {
		panic(err)
	}
	absPathCred := os.Getenv(key)

	log.Println(key, absPathCred)
}

//EKSClusterAccess ...
type EKSClusterAccess struct {
}

// Create ...
func (ca EKSClusterAccess) Create(createThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error) {
	defer clusterutil.TrackTime("Create EKS", time.Now())

	tagsCsv := clusterutil.ToCommaSeparateKeyValuePairs(createThis.Labels)
	err = eksctl.CreateClusterNoNodeGroup(createThis.Name, createThis.Location, createThis.K8sVersion, tagsCsv)
	if err != nil {
		err2 := ca.Delete(createThis)
		if err2 != nil {
			return nil, errors.Wrap(err2, "error deleting cluster after failing to create it; original error "+err.Error())
		}
		return nil, errors.Wrap(err, "cannot create cluster")

	}
	err = eksctl.AddLogging(createThis.Name, createThis.Location, createThis.K8sVersion, tagsCsv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot add logging")
	}
	for _, ng := range createThis.NodePools {
		err = eksctl.CreateNodeGroup(createThis.Name, ng.Name, createThis.Location, createThis.K8sVersion,
			ng.MachineType.Name, tagsCsv, ng.NodeCount, ng.DiskSizeGB, ng.Preemptible)
		if err != nil {
			err2 := ca.Delete(createThis)
			if err2 != nil {
				return nil, errors.Wrap(err2, "error deleting cluster after failing to create NodeGroup "+ng.Name+"; original error "+err.Error())
			}
			return nil, errors.Wrap(err, "cannot create NodeGroup "+ng.Name)
		}
	}
	created, err = ca.Describe(createThis)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe created cluster "+createThis.Name)
	}
	created.GeneratedBy = clusters.Created
	return created, nil
}

// Delete ...
func (ca EKSClusterAccess) Delete(deleteThis *clusters.ClusterInfo) (err error) {
	defer clusterutil.TrackTime("Delete EKS", time.Now())

	err = eksctl.DeleteCluster(deleteThis.Name, deleteThis.Location)
	if err != nil {
		return errors.Wrap(err, "cannot delete cluster")
	}
	return nil

}

//Describe ...
func (ca EKSClusterAccess) Describe(searchTemplate *clusters.ClusterInfo) (described *clusters.ClusterInfo, err error) {
	defer clusterutil.TrackTime("Describe EKS", time.Now())
	if searchTemplate.GeneratedBy == "" {
		searchTemplate.GeneratedBy = clusters.SearchTemplate
	}
	if searchTemplate.Location == "" {
		return nil, errors.New("must provide location to describe cluster")
	}
	clusterOutput, err := awssdk.DescribeCluster(searchTemplate.Name, searchTemplate.Location)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe cluster "+searchTemplate.Name)
	}
	described = clusterObjectToClusterInfo(clusterOutput, searchTemplate.Location, clusters.Read)
	describeNodegroupOutputs, err := awssdk.DescribeNodeGroups(searchTemplate.Name, searchTemplate.Location)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe nodes for cluster "+searchTemplate.Name)
	}
	if err := addNodeGroupObjectsAsNodePoolInfo(describeNodegroupOutputs, described); err != nil {
		return nil, errors.Wrap(err, "cannot add NodePoolInfos")
	}

	return described, nil
}

func addNodeGroupObjectsAsNodePoolInfo(eksNodeGroups []*eks.DescribeNodegroupOutput, cluster *clusters.ClusterInfo) error {
	for _, describeNodeGroup := range eksNodeGroups {
		nodeGroup := describeNodeGroup.Nodegroup
		if len(nodeGroup.InstanceTypes) > 1 {
			log.Println("NodeGroup has multiple InstanceTypes; support having only 1")

		}
		instanceType := *nodeGroup.InstanceTypes[0]
		machineType, err := GetMachineTypes().Get(instanceType) //TODO support multi-instance-type NG
		if err != nil {
			return errors.Wrap(err, "cannot get instance type "+instanceType)
		}
		if machineType.Name == "" {
			return errors.New("cannot find " + instanceType)
		}
		scaling := nodeGroup.ScalingConfig
		if *scaling.MaxSize != *scaling.DesiredSize || *scaling.MinSize != *scaling.DesiredSize {
			log.Println(fmt.Sprintf("Dynamic scaling unsupported: %v", scaling))
		}
		//NodeGroup name also available  nodeGroup.Labels["alpha.eksctl.io/nodegroup-name"]
		ngName := *nodeGroup.NodegroupName
		npi := clusters.NodePoolInfo{
			Name:        ngName,
			NodeCount:   int(*scaling.DesiredSize),
			K8sVersion:  *nodeGroup.Version,
			MachineType: machineType, //TODO Deal with missing instance types. What instance types are allowed in EKS?
			DiskSizeGB:  int(*nodeGroup.DiskSize),
		}

		cluster.AddNodePool(npi)
	}
	return nil
}

func clusterObjectToClusterInfo(clusterOutput *eks.DescribeClusterOutput, loc, generatedBy string) *clusters.ClusterInfo {
	cluster := clusterOutput.Cluster
	ci := &clusters.ClusterInfo{
		Scope:       "",
		Location:    loc,
		Name:        *cluster.Name,
		K8sVersion:  *cluster.Version,
		GeneratedBy: generatedBy,
		Cloud:       clusters.AWS,
		Labels:      clusterutil.StrPtrMapToStrMap(cluster.Tags),
	}

	return ci
}

// List ...
func (ca EKSClusterAccess) List(_, location string, tagFilter map[string]string) (listedClusters []*clusters.ClusterInfo, err error) {

	listedClusterNames, err := awssdk.DescribeClusters(location)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list clusters")
	}
	listedClusters = make([]*clusters.ClusterInfo, 0)
	unmatchedNames := make([]string, 0)
	matchedNames := make([]string, 0)
	for _, clusterOutput := range listedClusterNames {
		clusterName := *clusterOutput.Cluster.Name
		searchTemplate := &clusters.ClusterInfo{Cloud: clusters.AWS, Name: clusterName, Location: location, GeneratedBy: clusters.SearchTemplate}
		ci, err := ca.Describe(searchTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "cannot describe cluster "+clusterName)
		}

		match := clusterutil.LabelMatch(tagFilter, ci.Labels)
		if !match {
			log.Printf("Skipping cluster %s because labels do not match", ci.Name)
			unmatchedNames = append(unmatchedNames, ci.Name)
			continue
		}
		matchedNames = append(matchedNames, ci.Name)
		listedClusters = append(listedClusters, ci)
	}

	accessutil.PrintFilteringResults(clusters.AWS, tagFilter, matchedNames, unmatchedNames)
	return listedClusters, nil
}

// GetSupportedK8sVersions ...
func (ca EKSClusterAccess) GetSupportedK8sVersions(scope, location string) (versions []string, err error) {
	return []string{"1.12", "1.13", "1.14", "1.15"}, nil //TODO load dynamically

}

// MachineTypes ...
var eksMachineTypes *machinetypes.MachineTypes

// GetMachineTypes ...
func GetMachineTypes() *machinetypes.MachineTypes {
	return eksMachineTypes
}
func init() {
	var err error
	eksMachineTypes, err = loadMachineTypes()
	if err != nil {
		panic(fmt.Sprintf("cannot load EKS machine types %v", err))
	}
	if eksMachineTypes.Length() == 0 {
		panic(fmt.Sprintf("cannot load EKS machine types %v", err))
	}
}

func loadMachineTypes() (*machinetypes.MachineTypes, error) {
	ret := machinetypes.NewMachineTypeMap()

	fn := clusterutil.RootPath() + "/machine-types/aws-instance-types.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		wd, _ := os.Getwd()
		return nil, errors.Wrap(err, fmt.Sprintf("At %s: %v", wd, err))
	}

	r := csv.NewReader(csvfile)
	r.Comma = ','
	r.Comment = '#'
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
		ramGiBStr := record[2]
		ramGiBFloat, err := strconv.ParseFloat(ramGiBStr, 32)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse memory "+ramGiBStr)
		}

		cpus := record[3]
		cpuInteger, err := strconv.ParseInt(cpus, 10, 32)
		if err != nil || cpuInteger == 0 {
			return nil, errors.Wrap(err, "cannot parse cpus "+cpus)
		}
		ret.Set(name, machinetypes.MachineType{Name: name, CPU: int(cpuInteger), RAMMB: int(ramGiBFloat * 1024)})
	}
	return &ret, nil
}

//
