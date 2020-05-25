package access

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/eksctl"
	"clustercloner/clusters/util"
	"github.com/pkg/errors"
	"log"
	"os"
)

func init() {
	key := "AWS_SHARED_CREDENTIALS_FILE"
	cred := os.Getenv(key)
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
	for _, nodeGroup := range createThis.NodePools {
		err = eksctl.CreateNodeGroup(createThis.Name, nodeGroup.Name, createThis.Location, createThis.K8sVersion,
			nodeGroup.MachineType.Name, tagsCsv, int(nodeGroup.NodeCount), int(nodeGroup.DiskSizeGB))
		if err != nil {
			return nil, err
		}
	}
	created, err = ca.Describe(createThis)
	return created, err
}

// Delete ...
func (ca EKSClusterAccess) Delete(deleteThis *clusters.ClusterInfo) (err error) {
	log.Println("Implement") //TODO
	return
}

//Describe ...
func (ca EKSClusterAccess) Describe(describeThis *clusters.ClusterInfo) (created *clusters.ClusterInfo, err error) {

	log.Println("Implement") //TODO
	return
}

// List ...
func (ca EKSClusterAccess) List(subscription, location string, labelFilter map[string]string) (listedClusters []*clusters.ClusterInfo, err error) {
	log.Println("Implement") //TODO

	return
}

// GetSupportedK8sVersions ...
func (ca EKSClusterAccess) GetSupportedK8sVersions(scope, location string) (versions []string) {
	return []string{"1.12", "1.13", "1.14", "1.15"} //TODO load dynamically; handle the lack of minor version andpatch

}

// MachineTypeByName ...
func MachineTypeByName(machineType string) clusters.MachineType {
	return MachineTypes[machineType] //return zero object if not found
}

// MachineTypes ...
var MachineTypes map[string]clusters.MachineType
