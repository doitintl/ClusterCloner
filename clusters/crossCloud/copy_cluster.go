package crossCloud

import (
	"clusterCloner/clusters/eks"
	"clusterCloner/clusters/gke"

	_ "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	existing := gke.ReadClusters(cliCtx)
	eks.CreateCluster(existing.Clusters[0].Name)
}
