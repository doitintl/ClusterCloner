package access

import (
	containerv1 "cloud.google.com/go/container/apiv1"
	"clustercloner/clusters"
	"clustercloner/clusters/machinetypes"
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

// Delete ...
func (ca GKEClusterAccess) Delete(ci *clusters.ClusterInfo) error {
	bkgdCtx := context.Background()
	client, err := containerv1.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		return errors.Wrap(err, "cannot make client")
	}
	req := containerpb.DeleteClusterRequest{Name: projectLocationClusterPath(ci.Scope, ci.Location, ci.Name)}
	log.Println("About to delete GKE cluster", ci.Name)

	op, err := client.DeleteCluster(bkgdCtx, &req)
	if err != nil {
		return errors.Wrap(err, "cannot delete")
	}

	err = waitForClusterDeletion(ci.Scope, ci.Location, ci.Name, op.Name)
	if err != nil {
		return errors.Wrap(err, "waiting for cluster deletion")
	}
	return nil
}

// Describe ...
func (ca GKEClusterAccess) Describe(searchTemplate *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	if searchTemplate.GeneratedBy == "" {
		searchTemplate.GeneratedBy = clusters.SearchTemplate
	}
	if searchTemplate.GeneratedBy != clusters.SearchTemplate &&
		searchTemplate.GeneratedBy != clusters.Transformation { //In Create, we describe the created cluster based on the info used to create the cluster.
		panic(fmt.Sprintf("Wrong GeneratedBy: %s", searchTemplate.GeneratedBy))
	}
	cluster, err := getCluster(searchTemplate.Scope, searchTemplate.Location, searchTemplate.Name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster")
	}
	readClusterInfo, err := clusterObjectToClusterInfo(cluster, searchTemplate.Scope)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert cluster object")
	}
	readClusterInfo.SourceCluster = searchTemplate
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

// List lists clusters; location param can be region or zone
func (ca GKEClusterAccess) List(project, location string, labelFilter map[string]string) (ret []*clusters.ClusterInfo, err error) {
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

	unmatchedNames := make([]string, 0)
	matchedNames := make([]string, 0)
	for _, cluster := range resp.GetClusters() {

		match := clusterutil.LabelMatch(labelFilter, cluster.GetResourceLabels())
		if !match {
			unmatchedNames = append(unmatchedNames, cluster.GetName())
			continue
		}
		matchedNames = append(matchedNames, cluster.GetName())
		foundClusterInfo, err := clusterObjectToClusterInfo(cluster, project)
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert cluster object")
		}
		ret = append(ret, foundClusterInfo)
	}
	log.Printf("In listing GKE clusters, the label filter was %v. These matched %v; and these did not %v", labelFilter, matchedNames, unmatchedNames)

	return ret, nil

}

func clusterObjectToClusterInfo(clus *containerpb.Cluster, project string) (*clusters.ClusterInfo, error) {
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
		machineType, err := gkeMachineTypes.Get(np.GetConfig().MachineType)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get machien type "+np.GetConfig().MachineType)
		}
		if machineType.Name == "" { //zero-object
			return nil, errors.New("cannot find machine type " + np.GetConfig().MachineType)
		}
		npi := clusters.NodePoolInfo{
			Name:        np.GetName(),
			NodeCount:   int(np.GetInitialNodeCount()),
			MachineType: machineType,
			K8sVersion:  np.GetVersion(),
			DiskSizeGB:  int(np.GetConfig().GetDiskSizeGb()),
			Preemptible: np.GetConfig().Preemptible,
		}
		zero := machinetypes.MachineType{}
		if npi.MachineType == zero {
			panic("cannot read " + np.GetConfig().MachineType) //fix?
		}
		foundCluster.AddNodePool(npi)
	}
	return foundCluster, nil
}

func projectLocationPath(project, location string) string {
	return fmt.Sprintf("projects/%s/locations/%s", project, location)
}
func projectLocationClusterPath(project, location, clusterName string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, clusterName)
}
func projectLocationOperationPath(project, location, opName string) string {
	return fmt.Sprintf("projects/%s/locations/%s/operations/%s", project, location, opName)
}

// Create ...
func (ca GKEClusterAccess) Create(createThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {
	path := projectLocationPath(createThis.Scope, createThis.Location)

	var nodePools = make([]*containerpb.NodePool, len(createThis.NodePools))
	for i, npi := range createThis.NodePools {
		var nodeConfig = containerpb.NodeConfig{
			MachineType: npi.MachineType.Name,
			DiskSizeGb:  int32(npi.DiskSizeGB),
			Preemptible: npi.Preemptible,
		}
		np := containerpb.NodePool{
			Name:             npi.Name,
			Config:           &nodeConfig,
			InitialNodeCount: int32(npi.NodeCount),
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
	operation, err := clustMgrClient.CreateCluster(backgroundCtx, req)
	_ = operation // Could wait on the operation to be done instead of waiting on the cluster to be ready
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "cannot create")
	}

	createdCluster, err := ca.waitForClusterReadiness(createThis)
	if err != nil {
		return nil, errors.Wrap(err, "error in waiting for cluster to be ready")
	}
	return createdCluster, nil
}

func (ca GKEClusterAccess) waitForClusterReadiness(createThis *clusters.ClusterInfo) (*clusters.ClusterInfo, error) {

	var status containerpb.Cluster_Status
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
			return nil, errors.Errorf("Cluster in error status %s", status)
		default:
			panic(fmt.Sprintf("unknown status %s", status))
		}
	}
	createdCluster, err := ca.Describe(createThis) //redundant call to getCluster above,
	if err != nil {
		return nil, errors.Wrap(err, "could not describe cluster after creating it")
	}
	if createdCluster == nil {
		return nil, errors.New("createdCluster: nil")
	}
	createdCluster.GeneratedBy = clusters.Created
	createdCluster.SourceCluster = createThis
	return createdCluster, err
}
func getOperation(project, location, opName string) (*containerpb.Operation, error) {

	req := containerpb.GetOperationRequest{Name: projectLocationOperationPath(project, location, opName)}
	bkgdCtx := context.Background()
	client, err := containerv1.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot make client")
	}

	operation, err := client.GetOperation(bkgdCtx, &req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get the operation")
	}
	if operation == nil {
		return nil, errors.New("cannot get operation: nil")
	}

	return operation, nil
}
func waitForClusterDeletion(project, location, clusterName, opName string) error {
	var counter = -1
	log.Println("Waiting for deletion of " + clusterName + "; it may take a while")
	var status containerpb.Operation_Status
Waiting:
	for {
		time.Sleep(2 * time.Second)
		counter++
		op, err := getOperation(project, location, opName)
		if err != nil {
			return errors.Wrap(err, "cannot get operation to wait for shutdown of "+clusterName)
		}
		status = op.Status
		switch status {
		case containerpb.Operation_STATUS_UNSPECIFIED, containerpb.Operation_RUNNING, containerpb.Operation_PENDING:
			if counter%10 == 0 {
				log.Println("Waiting for deletion of "+clusterName+", operation status", status)
			}
			continue
		case containerpb.Operation_ABORTING:
			return errors.New("Aborting while waiting for shutdown")
		case containerpb.Operation_DONE:
			log.Printf("Deletion of GKE cluster is finished")
			break Waiting
		default:
			panic(fmt.Sprintf("unknown status %s", status))
		}
	}
	return nil
}
func init() {
	key := "GOOGLE_APPLICATION_CREDENTIALS"
	cred := os.Getenv(key)
	if cred == "" {
		log.Println(key + " not set; will use system gcloud authorization")
	} else {
		log.Println(key, "=", cred)
	}
}

// gkeMachineTypes ...
var gkeMachineTypes *machinetypes.MachineTypes

// GetMachineTypes ...
func GetMachineTypes() *machinetypes.MachineTypes {
	return gkeMachineTypes
}
func init() {
	var err error
	gkeMachineTypes, err = loadMachineTypes()
	if err != nil {
		panic(fmt.Sprintf("cannot load GKE machine types %v", err))
	}
	if gkeMachineTypes.Length() == 0 {
		panic(fmt.Sprintf("cannot load GKE machine types %v", err))
	}
}
func loadMachineTypes() (*machinetypes.MachineTypes, error) {
	ret := machinetypes.NewMachineTypeMap()
	fn := clusterutil.RootPath() + "/machine-types/gke-machine-types.csv"
	csvfile, err := os.Open(fn)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening "+fn)
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
		ramMB := ramGbFlt * 1000

		ret.Set(name, machinetypes.MachineType{Name: name, CPU: int(cpuInteger), RAMMB: int(ramMB)})
	}

	return &ret, nil
}

// supportedVersions ...
var supportedVersions []string

// GetSupportedK8sVersions ...
func (ca GKEClusterAccess) GetSupportedK8sVersions(scope, location string) (versions []string, err error) {

	if supportedVersions == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
		defer cancel()
		client, err := containerv1.NewClusterManagerClient(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create ClusterManagerClient")
		}

		supportedVersions = make([]string, 0)
		req := containerpb.GetServerConfigRequest{
			Name: projectLocationPath(scope, location),
		}
		resp, err := client.GetServerConfig(ctx, &req)
		if err != nil {
			return nil, errors.Wrap(err, "cannot GetServerConfig")
		}
		supportedVersions = resp.ValidMasterVersions[:]

	}
	return supportedVersions, nil
}
