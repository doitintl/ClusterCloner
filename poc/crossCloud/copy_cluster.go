package crossCloud

import (
	"clusterCloner/poc/aws"
	"clusterCloner/poc/google"

	_ "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	existing := google.ListClusters(cliCtx)
	aws.CreateCluster(existing.Clusters[0].Name)
}
