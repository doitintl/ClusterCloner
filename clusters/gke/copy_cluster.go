package gke

import (
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	origClusInfo := ReadClusters(cliCtx)
	CreateClusters(cliCtx, origClusInfo)

}
