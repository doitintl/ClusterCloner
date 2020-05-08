package access

import (
	container "cloud.google.com/go/container/apiv1"
	"clustercloner/clusters/clusterinfo"
	clusterutil "clustercloner/clusters/util"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"io"
	"log"
	"os"
	"strconv"
)

// GKEClusterAccess ...
type GKEClusterAccess struct {
}

// ListClusters lists clusters; location param can be region or zone
func (GKEClusterAccess) ListClusters(project, location string) (ret []clusterinfo.ClusterInfo, err error) {

	ret = make([]clusterinfo.ClusterInfo, 0)

	bkgdCtx := context.Background()
	client, err := container.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot make client")
	}

	path := fmt.Sprintf("projects/%s/locations/%s", project, location)
	req := &containerpb.ListClustersRequest{Parent: path}
	resp, err := client.ListClusters(bkgdCtx, req)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot list ")
	}

	for _, clus := range resp.GetClusters() {
		var nodePools = clus.GetNodePools()
		//var nodeCount int32 = 0
		foundCluster := clusterinfo.ClusterInfo{Scope: project,
			Location:            clus.Location,
			Name:                clus.Name,
			K8sVersion:          clus.CurrentMasterVersion,
			DeprecatedNodeCount: 1,
			GeneratedBy:         clusterinfo.READ,
			Cloud:               clusterinfo.GCP,
		}
		for _, np := range nodePools {

			nodePool := clusterinfo.NodePoolInfo{
				Name:        np.GetName(),
				NodeCount:   np.GetInitialNodeCount(),
				MachineType: ParseMachineType(np.GetConfig().MachineType),
				K8sVersion:  np.GetVersion(),
				DiskSizeGB:  np.GetConfig().GetDiskSizeGb(),
			}
			zero := clusterinfo.MachineType{}
			if nodePool.MachineType == zero {
				panic("cannot read " + np.GetConfig().MachineType)
			}
			foundCluster.AddNodePool(nodePool)
		}
		ret = append(ret, foundCluster)

	}
	return ret, nil

}

// ParseMachineType ...
func ParseMachineType(machineType string) clusterinfo.MachineType {
	ret := MachineTypes[machineType]
	return ret //return zero object if not found
}

// MachineTypes ...
var MachineTypes map[string]clusterinfo.MachineType

func init() {
	MachineTypes, _ = loadMachineTypes()

}
func loadMachineTypes() (map[string]clusterinfo.MachineType, error) {
	ret := make(map[string]clusterinfo.MachineType)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PWD", dir)
	fn := clusterutil.RootPath() + "/machine-types/gke-machine-types.csv"
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
		cpu := record[1]
		cpuFloat, err := strconv.ParseFloat(cpu, 32)
		if err != nil {
			return nil, err
		}
		cpuInt := int32(cpuFloat)

		ram := record[2]
		ramFlt, err := strconv.ParseFloat(ram, 32)
		if err != nil {
			return nil, err
		}
		ramInt := int32(ramFlt)

		ret[name] = clusterinfo.MachineType{Name: name, CPU: cpuInt, RAMGB: ramInt}
	}
	return ret, nil
} // CreateCluster ...

// CreateCluster ...
func (GKEClusterAccess) CreateCluster(createThis clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error) {

	initialNodeCount := createThis.DeprecatedNodeCount

	if initialNodeCount < 1 {
		log.Print("Copying a paused cluster, creating one node as a necessary minimum.")
		initialNodeCount = 1
	}
	path := fmt.Sprintf("projects/%s/locations/%s", createThis.Scope, createThis.Location)

	cluster := containerpb.Cluster{
		Name:                  createThis.Name,
		InitialNodeCount:      initialNodeCount,
		InitialClusterVersion: createThis.K8sVersion,
	}
	req := &containerpb.CreateClusterRequest{Parent: path, Cluster: &cluster}
	backgroundCtx := context.Background()
	clustMgrClient, _ := container.NewClusterManagerClient(backgroundCtx)
	resp, err := clustMgrClient.CreateCluster(backgroundCtx, req)
	if err != nil {
		log.Print(err)
		return clusterinfo.ClusterInfo{}, err
	}
	var created = createThis
	created.GeneratedBy = clusterinfo.CREATED
	log.Print(resp)
	return created, err
}
