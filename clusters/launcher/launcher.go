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
			Name: "inputfile",
			Usage: "Clusters to create in a JSON file. If this switch is used, other input switches (for reading from the cloud) cannot be used. " +
				"The file should represent an array of ClusterInfo, where these fields are mandatory in the JSON: " +
				"Cloud, Scope, Location, Name K8sVersion, NodePools. " +
				"These are optional: GeneratedBy, SourceCluster",
		},
		&cli.StringFlag{
			Name:  "inputscope",
			Usage: "GCP project or AKS resource group; for AWS value is ignored",
		},
		&cli.StringFlag{
			Name:  "outputscope",
			Usage: "GCP project or AKS resource group; for AWS value is ignored ",
		},
		&cli.StringFlag{
			Name:  "inputlocation",
			Usage: "GCP region (for regional clusters) or zone (zonal clusters); AWS region; or AKS region",
		},
		&cli.StringFlag{
			Name:  "inputcloud",
			Usage: "GCP, Azure, or AWS",
		},
		&cli.StringFlag{
			Name:  "outputcloud",
			Usage: "GCP, Azure, AWS, or Hub",
		},
		&cli.BoolFlag{
			Name:  "create",
			Usage: "true: Create new clusters; default is not to create (dry run)",
		},
		&cli.BoolFlag{
			Name:  "randomsuffix",
			Usage: "true: add a random suffix to cluster names to prevent collisions",
		},
		&cli.StringFlag{
			Name:  "labelfilter",
			Usage: "comma-separated list of key=value pairs; all must match for a clsuter to be included"},
	}
}

// Launch ...
func Launch(cliCtx *cli.Context) {
	googleCred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	log.Println("GOOGLE_APPLICATION_CREDENTIALS", googleCred)
	outputClusters, err := transformation.CloneFromCli(cliCtx)
	if err != nil {
		log.Fatalf("Error in transformation: %v", err)
	}

	outputString := util.ToJSON(outputClusters)
	exitCode, err := os.Stdout.WriteString(outputString + "\n")
	if err != nil {
		log.Fatalf("Error on exit %v, code %d", err, exitCode)
	}
}
