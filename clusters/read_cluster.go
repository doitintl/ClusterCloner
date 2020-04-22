package clusters

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
)

func ReadCluster(cliCtx *cli.Context) {

	proj := cliCtx.String("project")
	loc := cliCtx.String("location")
	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("projects/%s/locations/%s", proj, loc)
	req := &containerpb.ListClustersRequest{Parent: path}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	//	var clus = resp.Clusters
	var js []byte
	js, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println( string(js))
	}

}

