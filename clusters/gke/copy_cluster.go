package gke

import (
	"github.com/urfave/cli/v2"
)

func CopyCluster(cliCtx *cli.Context) {
	origClusInfo := ReadCluster(cliCtx)
	CreateClusters(cliCtx, origClusInfo)

}
