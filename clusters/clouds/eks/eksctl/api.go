package eksctl

import (
	"clustercloner/clusters"
	"clustercloner/clusters/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// CreateCluster ...
func CreateCluster(clusterName, region, k8sVersion, tagsCsv string) error {
	err := createClusterNoNodeGroup(clusterName, region, k8sVersion, tagsCsv)
	return err
}

func createClusterNoNodeGroup(clusterName, region, k8sVersion, tagsCsv string) error {
	args := []string{"eksctl", "create", "cluster", "--without-nodegroup", "--name", clusterName, "--region", region, "--tags", tagsCsv, "--version", k8sVersion}
	return runEksctl(args)
}

// AddLogging ...
func AddLogging(clusterName, region, k8sVersion, tagsCsv string) error {

	args := []string{"eksctl", "utils", "update-cluster-logging", "--cluster", clusterName, "--region", region}
	return runEksctl(args)
}

// CreateNodeGroup ...
func CreateNodeGroup(clusterName, nodeGroupName, region, k8sVersion, nodeInstanceType, tagsCsv string, nodeCount, diskSizeGB int) error {

	args := []string{"eksctl", "create", "nodegroup", "--managed",
		"--cluster", clusterName, "--name", nodeGroupName, "--region", region, "--version", k8sVersion,
		"--node-type", nodeInstanceType, "--tags", tagsCsv,
		"--nodes", strconv.Itoa(nodeCount), "--node-volume-size", strconv.Itoa(diskSizeGB),
	}
	return runEksctl(args)
}

// DeleteNodeGroup ...
func DeleteNodeGroup(clusterName, nodeGroupName, region string) error {
	//TODO is this even needed, or does DeleteCluster do it?
	return nil
}

// DeleteCluster ...
func DeleteCluster(clusterName, region string) error {
	args := []string{"eksctl", "delete", "cluster",
		"--name", clusterName, "--region", region,
		"--wait"}
	err := runEksctl(args)
	if err != nil {
		return errors.Wrap(err, "cannot runEksCtl for DeleteCluster")
	}
	return nil
}

// DescribeCluster ...
func DescribeCluster(clusterName, region string) (*clusters.ClusterInfo, error) {
	tempStdoutFile, oldStdout := util.ReplaceStdout()
	defer util.RestoreStdout(oldStdout)
	args := []string{"eksctl", "get", "cluster",
		"--name", clusterName, "--region", region, "--output", "json"}
	err := runEksctl(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for DescribeCluster")
	}
	httpResponse, err := ioutil.ReadFile(tempStdoutFile)

	if err != nil {
		return nil, errors.Wrap(err, "cannot load file with eksctl response "+tempStdoutFile)
	}
	if len(httpResponse) == 0 {
		return nil, errors.New("no output")
	}
	eksClus, err := parseClusterDescription(httpResponse)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse JSON in "+string(httpResponse))
	}
	_ = eksClus
	panic("implement by converting eksclus")
}

// ListClusters ...
func ListClusters(region, labelFilter string) ([]string, error) {
	tempStdoutFile, oldStdout := util.ReplaceStdout()
	defer util.RestoreStdout(oldStdout)
	args := []string{"eksctl", "get", "clusters", "--region", region, "--output", "table"}
	err := runEksctl(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for DescribeCluster")
	}
	outputAsTable, err := ioutil.ReadFile(tempStdoutFile)

	if err != nil {
		return nil, errors.Wrap(err, "cannot load input file "+tempStdoutFile)
	}
	if len(outputAsTable) == 0 {
		return nil, errors.New("no output")
	}
	err = os.Remove(tempStdoutFile)
	if err != nil {
		return nil, errors.Wrap(err, "cannot remove  "+tempStdoutFile)
	}
	outputAsTableStr := string(outputAsTable)

	log.Println(outputAsTableStr)
	clusterNames, err := parseClusterList(outputAsTableStr, region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse "+outputAsTableStr)
	}

	return clusterNames, nil
}

func resetOsArgs(oldOsArgs []string) {
	os.Args = oldOsArgs
}
func parseClusterList(s string, expectRegion string) ([]string, error) {
	ret := make([]string, 0)
	//example	"NAME\t\tREGION\nclus-sudic\tus-east-2\n"

	sNormalized := strings.ReplaceAll(s, "\t\t", "\t")
	lines := strings.Split(sNormalized, "\n")
	for i, line := range lines {
		parts := strings.Split(line, "\t")
		if line == "" {
			continue
		}
		if len(parts) != 2 {
			return nil, errors.New("wrong number of fields  " + line)
		}
		if i == 0 {
			if line != "NAME\tREGION" {
				return nil, errors.New("bad header line " + line)
			}
			continue
		}

		region := parts[1]
		if region != expectRegion {
			return nil, errors.New("unexpected region " + region + " instead of " + expectRegion)
		}
		clusterName := parts[0]
		ret = append(ret, clusterName)
	}

	return ret, nil
}
