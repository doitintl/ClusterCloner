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
			Usage:    "GCP region (for regional clusters) or zone (zonal clusters); AWS region; or AKS region",
			Required: true,
		},
		&cli.StringFlag{ //TODO allow inputting JSON for inputcloud=Hub
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
		&cli.BoolFlag{
			Name:  "randomsuffix",
			Usage: "true: add a random suffix to cluster names to prevent collisions",
		},
	}
}

// Launch ...
func Launch(cliCtx *cli.Context) {
	log.SetOutput(os.Stderr)
	outputClusters, err := transformation.Clone(cliCtx)
	if err != nil {
		log.Fatalf("Error in transformation %v", err)
	}
	var outputString string

	outputString = util.ToJSON(outputClusters)

	log.Println(outputString)

}
