package google

import (
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	origClusInfo := ListClusters(cliCtx)
	CreateClusters(cliCtx, origClusInfo)

}
