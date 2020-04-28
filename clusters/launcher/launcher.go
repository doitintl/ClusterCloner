package launcher

import (
	"clusterCloner/clusters/aks"
	"clusterCloner/clusters/util"
	"github.com/urfave/cli/v2"
)

func CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "scope",
			Usage:    "GCP project or AKS resource group",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "location",
			Usage: "GCP zone or AWS region",
		},
	}
}

func Launch(cliCtx *cli.Context) {
	//	azure.CreateClusterFromEnv("mycluster")
	proj := cliCtx.String("scope")
	_ = proj
	loc := cliCtx.String("location")
	//	azure.CreateClusterFromEnv("mycluster")
	ret, _ := aks.AksClusterAccess{}.ListClusters(proj, loc)
	util.PrintAsJson(ret)

}
