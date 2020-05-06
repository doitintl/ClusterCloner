package launcher

import (
	"clustercloner/clusters/transformation"
	"clustercloner/clusters/util"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// CLIFlags ...
func CLIFlags() []cli.Flag {
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
		&cli.StringFlag{
			Name:  "outputformat",
			Usage: "json or tf (Terraform). Default is json",
		},
	}
}

// Launch ...
func Launch(cliCtx *cli.Context) {
	log.SetOutput(os.Stderr)

	ret, _ := transformation.Clone(cliCtx)
	log.Println(util.MarshallToJSONString(ret))

}
