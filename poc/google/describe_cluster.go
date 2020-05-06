package google

import (
	container "cloud.google.com/go/container/apiv1"
	"clustercloner/poc/utilities"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
)

//ListClusters Return data on the cluster in JSON form. cliCtx shold provide project and location (zone, where _ means all zones)
func ListClusters(cliCtx *cli.Context) *containerpb.ListClustersResponse {

	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	proj := cliCtx.String("project")
	loc := cliCtx.String("location")
	path := fmt.Sprintf("projects/%s/locations/%s", proj, loc)
	req := &containerpb.ListClustersRequest{Parent: path}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	utilities.PrintAsJSON(resp)
	return resp
}
