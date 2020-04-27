package gke

import (
	containers "cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
)

//CreateClusters Create a cluster with the given params
func CreateClusters(cliCtx *cli.Context, origClustersInfo *containerpb.ListClustersResponse) {
	//todo support Azure, AWS, for both read and write cluster
	ctx := context.Background()
	clustMgrClient, err := containers.NewClusterManagerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	proj := cliCtx.String("project")
	loc := cliCtx.String("location")

	if loc == "=" {
		panic("Cannot use wildcard for zones (_) to create cluster")
	}
	for _, clusterInfo := range origClustersInfo.Clusters {
		createCluster(ctx, proj, loc, clusterInfo, clustMgrClient)
	}

}

func createCluster(bkgrdCtx context.Context, proj string, loc string, origCluster *containerpb.Cluster, clustMgrClient *containers.ClusterManagerClient) {
	clusterName := origCluster.Name + "-copy"
	initialNodeCount := origCluster.InitialNodeCount
	if initialNodeCount < 1 {
		log.Print("Copying a paused cluster, creating one node as a necessary minimum.")
		initialNodeCount = 1
	}
	path := fmt.Sprintf("projects/%s/locations/%s", proj, loc)

	cluster := containerpb.Cluster{
		Name:             clusterName,
		InitialNodeCount: initialNodeCount,
	}
	req := &containerpb.CreateClusterRequest{Parent: path, Cluster: &cluster}
	resp, err := clustMgrClient.CreateCluster(bkgrdCtx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(resp)
}
