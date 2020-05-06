package crosscloud

import (
	"clustercloner/poc/aws"
	"clustercloner/poc/azure"
	"github.com/urfave/cli/v2"
)

// CliFlags ...
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

// PocLaunch ...
func PocLaunch() {
	azure.CreateClusterFromEnv("mycluster2")
	//_, _ = azure.describeCluster("joshua-playground", "mycluster")

	aws.DescribeNG()
	//aws.describeCluster("mycluster")
	//aws.CreateCluster("cluster3")

}
