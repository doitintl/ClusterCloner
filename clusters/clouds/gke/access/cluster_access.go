package access

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

// DescribeCluster ...
func (ca GKEClusterAccess) DescribeCluster(describeThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	if describeThis.GeneratedBy == "" {
		describeThis.GeneratedBy = clusters.SearchTemplate
	}
	if describeThis.GeneratedBy != clusters.SearchTemplate &&
		describeThis.GeneratedBy != clusters.Transformation { //In CreateCluster, we describe the created cluster based on the info used to create the cluster.
		panic(fmt.Sprintf("Wrong GeneratedBy: %s", describeThis.GeneratedBy))
	}
	cluster, err := getCluster(describeThis.Scope, describeThis.Location, describeThis.Name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster")
	}
	readClusterInfo := clusterObjectToClusterInfo(cluster, describeThis.Scope)
	readClusterInfo.SourceCluster = describeThis
	return readClusterInfo, nil
}

func getCluster(project, location, name string) (*containerpb.Cluster, error) {
	req := containerpb.GetClusterRequest{
		Name: projectLocationClusterPath(project, location, name),
	}
	bkgdCtx := context.Background()
	client, err := containerv1.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot make client")
	}

	cluster, err := client.GetCluster(bkgdCtx, &req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster")
	}
	return cluster, nil
}

// ListClusters lists clusters; location param can be region or zone
func (ca GKEClusterAccess) ListClusters(project, location string, labelFilter map[string]string) (ret []*clusters.ClusterInfo, err error) {
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

	for _, cluster := range resp.GetClusters() {

		match := clusterutil.LabelMatch(labelFilter, cluster.GetResourceLabels())
		if !match {
			log.Printf("Skipping cluster %s because labels do not match", cluster.GetName())
			continue
		}
		foundClusterInfo := clusterObjectToClusterInfo(cluster, project)
		ret = append(ret, foundClusterInfo)
	}
	return ret, nil

}

func clusterObjectToClusterInfo(clus *containerpb.Cluster, project string) *clusters.ClusterInfo {
	labels := clus.GetResourceLabels()
	foundCluster := &clusters.ClusterInfo{
		Scope:       project,
		Location:    clus.Location,
		Name:        clus.Name,
		K8sVersion:  clus.CurrentMasterVersion,
		GeneratedBy: clusters.Read,
		Labels:      clusterutil.CopyStringMap(labels),
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
			panic("cannot read " + np.GetConfig().MachineType) //fix?
		}
		foundCluster.AddNodePool(nodePool)
	}
	return foundCluster
}

func projectLocationPath(project, location string) string {
	path := fmt.Sprintf("projects/%s/locations/%s", project, location)
	return path
}
func projectLocationClusterPath(project, location, clusterName string) string {
	path := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, clusterName)
	return path
}

// CreateCluster ...
func (ca GKEClusterAccess) CreateCluster(createThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
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
		ResourceLabels:        createThis.Labels,
	}
	req := &containerpb.CreateClusterRequest{Parent: path, Cluster: &cluster}

	backgroundCtx := context.Background()
	clustMgrClient, _ := containerv1.NewClusterManagerClient(backgroundCtx)
	resp, err := clustMgrClient.CreateCluster(backgroundCtx, req)
	_ = resp
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot create")
	}

	createdCluster, err2 := ca.waitForClusterReadiness(createThis)
	if err2 != nil {
		return nil, errors.Wrap(err2, "error in waiting for cluster to be ready")
	}
	return createdCluster, err
}

func (ca GKEClusterAccess) waitForClusterReadiness(createThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	var status = containerpb.Cluster_STATUS_UNSPECIFIED
	var err error
Waiting:
	for {
		time.Sleep(time.Second)
		clusterPb, err := getCluster(createThis.Scope, createThis.Location, createThis.Name)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get created cluster")
		}
		status = clusterPb.Status
		switch status {
		case containerpb.Cluster_STATUS_UNSPECIFIED, containerpb.Cluster_PROVISIONING, containerpb.Cluster_RECONCILING:
			continue
		case containerpb.Cluster_RUNNING:
			log.Printf("Cluster %s now running", createThis.Name)
			break Waiting
		case containerpb.Cluster_ERROR, containerpb.Cluster_STOPPING, containerpb.Cluster_DEGRADED:
			return nil, errors.New(fmt.Sprintf("Cluster in error status %s", status))
		default:
			panic(fmt.Sprintf("unknown status %s", status))
		}
	}
	createdCluster, err := ca.DescribeCluster(createThis) //redundant call to getCluster above,
	if err != nil {
		return nil, errors.Wrap(err, "could not describe cluster after creating it")
	}
	if createdCluster == nil {
		return nil, errors.New("createdCluster nil")
	}
	createdCluster.GeneratedBy = clusters.Created
	createdCluster.SourceCluster = createThis
	return createdCluster, err
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
	if MachineTypes == nil || len(MachineTypes) == 0 || err != nil {
		panic(fmt.Sprintf("cannot load machine types %v", err))
	}
}
func loadMachineTypes() (map[string]clusters.MachineType, error) {
	ret := make(map[string]clusters.MachineType)
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
		cpus := record[1]
		cpuInteger, err := strconv.ParseInt(cpus, 10, 32)
		if err != nil || cpuInteger == 0 {
			return nil, errors.Wrap(err, "cannot parse cpus "+cpus)
		}
		ramGBStr := record[2]
		ramGbFlt, err := strconv.ParseFloat(ramGBStr, 32)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse "+ramGBStr)
		}
		ramMB := int32(ramGbFlt * 1000)

		ret[name] = clusters.MachineType{Name: name, CPU: int32(cpuInteger), RAMMB: ramMB}
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
		supportedVersions = resp.ValidMasterVersions[:] //TODO use .ValidNodeVersionsto supply versions to nodes

	}
	return supportedVersions
}
