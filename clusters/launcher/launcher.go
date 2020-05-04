package launcher

import (
	"clusterCloner/clusters/cluster_transformation"
	"clusterCloner/clusters/util"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "inputscope",
			Usage: "GCP project or AKS resource group; for AWS value is ignored",
		},
		&cli.StringFlag{
			Name:     "outputscope",
			Usage:    "GCP project or AKS resource group; for AWS value is ignored ",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "inputlocation",
			Usage:    "GCP zone or AWS region or AKS region",
			Required: true,
		},
		&cli.StringFlag{ //todo allow inputting JSON for inputcloud=Hub
			Name:     "inputcloud",
			Usage:    "GCP, Azure, or AWS",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "outputcloud",
			Usage:    "GCP, Azure, AWS, or Hub",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "create",
			Usage: "true: Create new clusters; default is not to create (dry run)",
		},
	}
}

func Launch(cliCtx *cli.Context) {
	log.SetOutput(os.Stderr)

	ret, _ := cluster_transformation.Clone(cliCtx)
	//	ret, _ := access.AksClusterAccess{}.ListClusters(scope, loc)
	log.Println(util.MarshallToJsonString(ret))

}
