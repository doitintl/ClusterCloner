package access

import (
	container "cloud.google.com/go/container/apiv1"
	"clustercloner/clusters/clusterinfo"
	"context"
	"fmt"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
)

// GkeClusterAccess ...
type GkeClusterAccess struct {
}

// ListClusters lists clusters; location param can be region or zone
func (GkeClusterAccess) ListClusters(project, location string) (ret []clusterinfo.ClusterInfo, err error) {

	ret = make([]clusterinfo.ClusterInfo, 0)

	bkgdCtx := context.Background()
	client, err := container.NewClusterManagerClient(bkgdCtx)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("projects/%s/locations/%s", project, location)
	req := &containerpb.ListClustersRequest{Parent: path}
	resp, err := client.ListClusters(bkgdCtx, req)
	if err != nil {
		log.Fatal(err)
	}

	for _, clus := range resp.GetClusters() {
		var nodePools = clus.GetNodePools()
		var nodeCount int32 = 0
		for _, np := range nodePools {
			nodeCount += np.InitialNodeCount
		}

		foundCluster := clusterinfo.ClusterInfo{Scope: project,
			Location:    location,
			Name:        clus.Name,
			NodeCount:   nodeCount,
			K8sVersion:  clus.CurrentMasterVersion,
			GeneratedBy: clusterinfo.READ,
			Cloud:       clusterinfo.GCP,
		}
		ret = append(ret, foundCluster)

	}
	return ret, nil

}

// CreateCluster ...
func (GkeClusterAccess) CreateCluster(createThis clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error) {

	initialNodeCount := createThis.NodeCount

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
