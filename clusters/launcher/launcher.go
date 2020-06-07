package launcher

import (
	"clustercloner/clusters/transformation"
	"clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
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
				"These fields of ClusterInfo are optional and not used: GeneratedBy, SourceCluster; as well as these fields of MachineType: RAMMB and CPU",
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
			Name:  "nodryrun",
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
func Launch(cliCtx *cli.Context) (err error) {

	outputClusters, err := transformation.CloneFromCli(cliCtx)
	if err != nil {
		return errors.Wrap(err, "error in transformation")
	}

	outputString := util.ToJSON(outputClusters)
	errorCode, err := os.Stdout.WriteString(outputString + "\n")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error writing output (error code %d)", errorCode))
	}
	return nil
}
