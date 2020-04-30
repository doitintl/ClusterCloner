package launcher

import (
	"clusterCloner/clusters/clouds/aks/access"
	"clusterCloner/clusters/util"
	"github.com/urfave/cli/v2"
)

func CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "inputscope",
			Usage:    "GCP project or AKS resource group",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "outputscope",
			Usage:    "GCP project or AKS resource group",
			Required: false,
		},
		&cli.StringFlag{
			Name:  "inputlocation",
			Usage: "GCP zone or AWS region or AKS region",
		},
		&cli.StringFlag{
			Name:     "inputcloud",
			Usage:    "GCP, Azure, or AWS",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "outputcloud",
			Usage:    "GCP, Azure, or AWS",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "true for JSON output only",
		},
	}
}

func Launch(cliCtx *cli.Context) {
	scope := cliCtx.String("scope")
	loc := cliCtx.String("location")
	ret, _ := access.AksClusterAccess{}.ListClusters(scope, loc)
	util.PrintAsJson(ret)

}
