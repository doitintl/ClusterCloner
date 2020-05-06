package google

import (
	"github.com/urfave/cli/v2"
)

// CopyCluster ...
func CopyCluster(cliCtx *cli.Context) {
	origClusInfo := ListClusters(cliCtx)
	CreateClusters(cliCtx, origClusInfo)

}
