package accessgke

import (
	containerv1 "cloud.google.com/go/container/apiv1"
	"clustercloner/clusters"
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
	"time"
)

// GKEClusterAccess ...
type GKEClusterAccess struct {
}

// ListClusters lists clusters; location param can be region or zone
func (ca GKEClusterAccess) ListClusters(project, location string) (ret []*clusters.ClusterInfo, err error) {

	ret = make([]*clusters.ClusterInfo, 0)

	bkgdCtx := context.Background()
	client, err := containerv1.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot make client")
	}

	path := projectLocationPath(project, location)
	req := &containerpb.ListClustersRequest{Parent: path}
	resp, err := client.ListClusters(bkgdCtx, req)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot list")
	}

	for _, clus := range resp.GetClusters() {

		foundCluster := &clusters.ClusterInfo{Scope: project,
			Location:    clus.Location,
			Name:        clus.Name,
			K8sVersion:  clus.CurrentMasterVersion,
			GeneratedBy: clusters.READ,
			Cloud:       clusters.GCP,
		}

		var nodePools = clus.GetNodePools()
		for _, np := range nodePools {
			nodePool := clusters.NodePoolInfo{
				Name:        np.GetName(),
				NodeCount:   np.GetInitialNodeCount(),
				MachineType: MachineTypeByName(np.GetConfig().MachineType),
				K8sVersion:  np.GetVersion(),
				DiskSizeGB:  np.GetConfig().GetDiskSizeGb(),
			}
			zero := clusters.MachineType{}
			if nodePool.MachineType == zero {
				panic("cannot read " + np.GetConfig().MachineType)
			}
			foundCluster.AddNodePool(nodePool)
		}
		ret = append(ret, foundCluster)

	}
	return ret, nil

}

func projectLocationPath(project string, location string) string {
	path := fmt.Sprintf("projects/%s/locations/%s", project, location)
	return path
}

// CreateCluster ...
func (GKEClusterAccess) CreateCluster(createThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	path := projectLocationPath(createThis.Scope, createThis.Location)

	var nodePools = make([]*containerpb.NodePool, len(createThis.NodePools))
	for i, npi := range createThis.NodePools {
		var nodeConfig = containerpb.NodeConfig{
			MachineType: npi.MachineType.Name,
			DiskSizeGb:  npi.DiskSizeGB,
		}
		np := containerpb.NodePool{
			Name:             npi.Name,
			Config:           &nodeConfig,
			InitialNodeCount: npi.NodeCount,
			Version:          npi.K8sVersion,
		}
		nodePools[i] = &np
	}
	cluster := containerpb.Cluster{
		Name:                  createThis.Name,
		InitialClusterVersion: createThis.K8sVersion,
		NodePools:             nodePools,
	}
	req := &containerpb.CreateClusterRequest{Parent: path, Cluster: &cluster}

	backgroundCtx := context.Background()
	clustMgrClient, _ := containerv1.NewClusterManagerClient(backgroundCtx)
	resp, err := clustMgrClient.CreateCluster(backgroundCtx, req)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot create")
	}
	//todo check status in a loops so that this is synchronous
	createThis.GeneratedBy = clusters.CREATED
	log.Println(resp)
	return createThis, err
}

// MachineTypeByName ... //todo inline this and other MacineTypeByName
func MachineTypeByName(machineType string) clusters.MachineType {
	return MachineTypes[machineType] //return zero object if not found
}

// MachineTypes ...
var MachineTypes map[string]clusters.MachineType

func init() {
	MachineTypes, _ = loadMachineTypes()

}
func loadMachineTypes() (map[string]clusters.MachineType, error) {
	ret := make(map[string]clusters.MachineType)
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

		ret[name] = clusters.MachineType{Name: name, CPU: cpuInt, RAMGB: ramInt}
	}
	return ret, nil
}

// supportedVersions ...
var supportedVersions []string

// GetSupportedK8sVersions ...
func (ca GKEClusterAccess) GetSupportedK8sVersions(scope, location string) []string {

	if supportedVersions == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
		defer cancel()
		client, err := containerv1.NewClusterManagerClient(ctx)
		if err != nil {
			log.Println(err)
			return nil
		}

		supportedVersions = make([]string, 0)
		req := containerpb.GetServerConfigRequest{
			Name: projectLocationPath(scope, location),
		}
		resp, err := client.GetServerConfig(ctx, &req)
		if err != nil {
			log.Println(err)
			return nil
		}
		supportedVersions = resp.ValidMasterVersions[:] //todo use .ValidNodeVersionsto supply versions to nodes

	}
	return supportedVersions
}
