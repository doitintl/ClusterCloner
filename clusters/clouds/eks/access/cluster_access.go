package access

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/eksctl"
	"clustercloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	key := "AWS_SHARED_CREDENTIALS_FILE"
	cred := os.Getenv(key)
	if cred == "" {
		log.Println("No " + key + " env variable, so using awscredentials in application root as default")
		cred = "awscredentials"
	}
	rootPath := util.RootPath() + "/" + cred
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
	tagsCsv := util.ToCommaSeparateKeyValuePairs(createThis.Labels)
	err = eksctl.CreateCluster(createThis.Name, createThis.Location, createThis.K8sVersion, tagsCsv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create cluster")
	}
	err = eksctl.AddLogging(createThis.Name, createThis.Location, createThis.K8sVersion, tagsCsv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot add logging")
	}
	for _, ng := range createThis.NodePools {
		//TODO support spot instances. Not now available in Managed Node Groups. See https: //github.com/aws/containers-roadmap/issues/583
		err = eksctl.CreateNodeGroup(createThis.Name, ng.Name, createThis.Location, createThis.K8sVersion,
			ng.MachineType.Name, tagsCsv, int(ng.NodeCount), int(ng.DiskSizeGB), ng.Preemptible)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create NodeGroup "+ng.Name)
		}
	}
	created, err = ca.Describe(createThis)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe created cluster "+createThis.Name)
	}
	created.GeneratedBy = clusters.Created
	return created, err
}

// Delete ...
func (ca EKSClusterAccess) Delete(deleteThis *clusters.ClusterInfo) (err error) {
	//TODO maybe delete NodeGroups separately
	//TODO check if VPC gets deleted
	err = eksctl.DeleteCluster(deleteThis.Name, deleteThis.Location)
	if err != nil {
		return errors.Wrap(err, "cannot delete cluster")
	}
	return nil

}

//Describe ...
func (ca EKSClusterAccess) Describe(searchTemplate *clusters.ClusterInfo) (described *clusters.ClusterInfo, err error) {
	if searchTemplate.GeneratedBy == "" {
		searchTemplate.GeneratedBy = clusters.SearchTemplate
	}
	if searchTemplate.Location == "" {
		return nil, errors.New("must provide location to describe cluster")
	}
	eksCluster, err := eksctl.DescribeCluster(searchTemplate.Name, searchTemplate.Location)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe cluster "+searchTemplate.Name)
	}
	described = clusterObjectToClusterInfo(eksCluster, searchTemplate.Location, clusters.Read)
	eksNodes, err := eksctl.DescribeNodeGroups(searchTemplate.Name, searchTemplate.Location)
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe nodes for cluster "+searchTemplate.Name)
	}
	if err := addNodeGroupObjectsAsNodePoolInfo(eksNodes, described); err != nil {
		return nil, errors.New("cannot add NodePoolInfos")
	}

	return described, nil
}

func addNodeGroupObjectsAsNodePoolInfo(eksNodeGroups []eksctl.EKSNodeGroup, described *clusters.ClusterInfo) error {
	for _, eksNg := range eksNodeGroups {
		mt := MachineTypeByName(eksNg.InstanceType)
		if mt.Name == "" {
			return errors.New("cannot find " + eksNg.InstanceType)
		}
		npi := clusters.NodePoolInfo{
			Name:        eksNg.Name,
			NodeCount:   int(eksNg.DesiredCapacity),
			K8sVersion:  described.K8sVersion, //TODO: Is this available  per-NodeGroup? It is not in eksctl output
			MachineType: mt,                   //TODO deal with missing instance tyoe ereT
			DiskSizeGB:  22,                   //TODO Find this data. How? It is not in eksctl output
		}

		described.AddNodePool(npi)
	}
	return nil
}

func clusterObjectToClusterInfo(eksClus eksctl.EKSCluster, loc, generatedBy string) *clusters.ClusterInfo {
	ci := &clusters.ClusterInfo{
		Scope:       "",
		Location:    loc,
		Name:        eksClus.Name,
		K8sVersion:  eksClus.Version,
		GeneratedBy: generatedBy,
		Cloud:       clusters.Azure,
		Labels:      eksClus.Tags,
	}

	return ci
}

// List ...
func (ca EKSClusterAccess) List(_, location string, tagFilter map[string]string) (listedClusters []*clusters.ClusterInfo, err error) {

	tagsCsv := util.ToCommaSeparateKeyValuePairs(tagFilter)
	listedClusterNames, err := eksctl.ListClusters(location, tagsCsv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list clusters")
	}
	listedClusters = make([]*clusters.ClusterInfo, 0)
	unmatchedNames := make([]string, 0)
	matchedNames := make([]string, 0)
	for _, clusterName := range listedClusterNames {
		searchTemplate := &clusters.ClusterInfo{Cloud: clusters.AWS, Name: clusterName, Location: location, GeneratedBy: clusters.SearchTemplate}
		ci, err := ca.Describe(searchTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "cannot describe cluster "+clusterName)
		}

		match := util.LabelMatch(tagFilter, ci.Labels)
		if !match {
			log.Printf("Skipping cluster %s because labels do not match", ci.Name)
			unmatchedNames = append(unmatchedNames, ci.Name)
			continue
		}
		matchedNames = append(matchedNames, ci.Name)

		listedClusters = append(listedClusters, ci)
		//TODO describe nodegroups too?
	}

	log.Printf("In listing clusters, these matched the label filter %v; and these did not %v\n", matchedNames, unmatchedNames)

	return listedClusters, nil
}

// GetSupportedK8sVersions ...
func (ca EKSClusterAccess) GetSupportedK8sVersions(scope, location string) (versions []string) {
	return []string{"1.12", "1.13", "1.14", "1.15"} //TODO load dynamically; handle the lack of minor version andpatch

}

// MachineTypeByName ...
func MachineTypeByName(machineType string) clusters.MachineType {
	//TODO load EKS Machine Types at init
	return MachineTypes[machineType]
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

	fn := util.RootPath() + "/machine-types/aws-instance-types.csv"
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
		// API Name; Display Name; Memory GiB; vCPUs; Supports EKS
		supportsEks := record[4]
		if strings.ToLower(supportsEks) != "true" {
			continue

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
		ret[name] = clusters.MachineType{Name: name, CPU: int(cpuInteger), RAMMB: int(ramGiBFloat * 1024)}
	}
	return ret, nil
}

//
