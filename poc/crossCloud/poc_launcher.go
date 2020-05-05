package crossCloud

import (
	"clusterCloner/poc/aws"
	"clusterCloner/poc/azure"
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

func PocLaunch() {
	azure.CreateClusterFromEnv("mycluster2")
	//_, _ = azure.describeCluster("joshua-playground", "mycluster")

	aws.DescribeNG()
	//aws.describeCluster("mycluster")
	//aws.CreateCluster("cluster3")

}
