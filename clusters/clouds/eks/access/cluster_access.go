package access

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/awssdk"
	"clustercloner/clusters/clouds/eks/eksctl"
	"clustercloner/clusters/machinetypes"
	"clustercloner/clusters/util"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strconv"
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
		//TODO support spot instances. Not now available in Managed Node Groups. See https: //github.com/aws/containers-roadmap/issues/583
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
	for _, describeNg := range eksNodeGroups {
		ng := describeNg.Nodegroup
		if len(ng.InstanceTypes) > 1 {
			log.Println("NodeGroup has multiple InstanceTypes; support having only 1")

		}
		mt := MachineTypeByName(*ng.InstanceTypes[0]) //TODO support multi-instance-type NG
		if mt.Name == "" {
			return errors.New("cannot find " + *ng.InstanceTypes[0])
		}
		scaling := ng.ScalingConfig
		if *scaling.MaxSize != *scaling.DesiredSize || *scaling.MinSize != *scaling.DesiredSize {
			log.Println(fmt.Sprintf("Dynamic scaling unsupported: %v", scaling))
		}
		ngName := *ng.Labels["alpha.eksctl.io/nodegroup-name"]

		npi := clusters.NodePoolInfo{
			Name:        ngName,
			NodeCount:   int(*scaling.DesiredSize),
			K8sVersion:  *ng.Version, //TODO: Is this available  per-NodeGroup? It is not in eksctl output
			MachineType: mt,          //TODO Deal with missing instance types
			DiskSizeGB:  int(*ng.DiskSize),
		}

		cluster.AddNodePool(npi)
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
		Cloud:       clusters.AWS,
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
	}

	log.Printf("In listing clusters, these matched the label filter %v; and these did not %v\n", matchedNames, unmatchedNames)

	return listedClusters, nil
}

// GetSupportedK8sVersions ...
func (ca EKSClusterAccess) GetSupportedK8sVersions(scope, location string) (versions []string, err error) {
	return []string{"1.12", "1.13", "1.14", "1.15"}, nil //TODO load dynamically; handle the lack of minor version andpatch

}

// MachineTypeByName ...
// TODO Reduce is repetition of the MachineTYpes code in EKS/AKS/GKE
func MachineTypeByName(machineType string) machinetypes.MachineType {
	mt, err := MachineTypes.Get(machineType)
	if err != nil {
		log.Println("cannot get " + machineType + "; " + err.Error())
		return machinetypes.MachineType{}
	}
	return mt
}

// MachineTypes ...
var MachineTypes *machinetypes.MachineTypeMap

func init() {
	var err error
	MachineTypes, err = loadMachineTypes()

	if err != nil && MachineTypes.Length() == 0 {
		panic(fmt.Sprintf("cannot load machine types %v", err))
	}
}

func loadMachineTypes() (*machinetypes.MachineTypeMap, error) {
	ret := machinetypes.NewMachineTypeMap()

	fn := util.RootPath() + "/machine-types/aws-instance-types.csv"
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
