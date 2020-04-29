package gke

import (
	container "cloud.google.com/go/container/apiv1"
	"clusterCloner/clusters"
	"context"
	"fmt"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
)

type GkeClusterAccess struct {
}

func (GkeClusterAccess) ListClusters(project, location string) (ret []clusters.ClusterInfo, err error) {
	ret = make([]clusters.ClusterInfo, 0)

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
		clusInfo := clusters.ClusterInfo{Scope: project, Location: location, Name: clus.Name, NodeCount: clus.InitialNodeCount}
		ret = append(ret, clusInfo)

	}
	return ret, nil
}
func (GkeClusterAccess) CreateCluster(clusterInfo clusters.ClusterInfo) error {

	initialNodeCount := clusterInfo.NodeCount

	if initialNodeCount < 1 {
		log.Print("Copying a paused cluster, creating one node as a necessary minimum.")
		initialNodeCount = 1
	}
	path := fmt.Sprintf("projects/%s/locations/%s", clusterInfo.Scope, clusterInfo.Location)

	cluster := containerpb.Cluster{
		Name:             clusterInfo.Name,
		InitialNodeCount: initialNodeCount,
	}
	req := &containerpb.CreateClusterRequest{Parent: path, Cluster: &cluster}
	backgroundCtx := context.Background()
	clustMgrClient, _ := container.NewClusterManagerClient(backgroundCtx)
	resp, err := clustMgrClient.CreateCluster(backgroundCtx, req)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(resp)
	return nil
}
