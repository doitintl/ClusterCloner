package crossCloud

import (
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
	//azure.CreateClusterFromEnv("mycluster")
	_, _ = azure.DescribeCluster("joshua-playground", "mycluster")
	//aws.DescribeCluster("mycluster")
	//aws.CreateCluster("cluster3")

}
