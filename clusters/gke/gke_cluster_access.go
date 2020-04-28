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
		clusInfo := clusters.ClusterInfo{clus.Name, clus.InitialNodeCount}
		ret = append(ret, clusInfo)

	}
	return ret, nil
}
