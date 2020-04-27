package crossCloud

import (
	"clusterCloner/clusters/aks"
	"clusterCloner/clusters/gke"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	aksCluster, _ := aks.ReadCluster("joshua-playground", "mycluster")
	mcp := aksCluster.ManagedClusterProperties
	var app *[]containerservice.AgentPoolProfile = mcp.AgentPoolProfiles
	agentPool := (*app)[0]
	nodeCount := agentPool.Count
	gke.CreateCluster(cliCtx, "mycluster", *nodeCount)
}
