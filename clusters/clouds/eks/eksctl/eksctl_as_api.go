package eksctl

import (
	"clustercloner/clusters/util"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// CreateClusterNoNodeGroup ...
func CreateClusterNoNodeGroup(clusterName, region, k8sVersion, tagsCsv string) error {
	args := []string{"eksctl", "create", "cluster", "--without-nodegroup", "--name", clusterName, "--region", region, "--tags", tagsCsv, "--version", k8sVersion}
	return runEksctl(args)
}

// AddLogging ...
func AddLogging(clusterName, region, k8sVersion, tagsCsv string) error {
	args := []string{"eksctl", "utils", "update-cluster-logging", "--cluster", clusterName, "--region", region, "--enable-types", "all"}
	return runEksctl(args)
}

// CreateNodeGroup ...
func CreateNodeGroup(clusterName, nodeGroupName, region, k8sVersion, nodeInstanceType, tagsCsv string, nodeCount, diskSizeGB int, preemptible bool) error {

	args := []string{"eksctl", "create", "nodegroup", "--managed",
		"--cluster", clusterName, "--name", nodeGroupName, "--region", region, "--version", k8sVersion,
		"--node-type", nodeInstanceType, "--tags", tagsCsv,
		"--nodes", strconv.Itoa(nodeCount), "--node-volume-size", strconv.Itoa(diskSizeGB),
	}
	err := runEksctl(args)
	if err != nil {
		return errors.Wrap(err, "cannot runEksCtl for CreateNodeGroup")
	}
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

// DescribeCluster ...TODO Replace
func DescribeCluster(clusterName, region string) (EKSCluster, error) {
	response, err := runDescribeCluster(clusterName, region)
	zero := EKSCluster{}
	if err != nil {
		return zero, errors.Wrap(err, "cannot runDescribeCluster")
	}
	if len(response) == 0 {
		return zero, errors.New("no output")
	}
	eksClusters, err := parseClusterDescription(response)
	if err != nil {
		return zero, errors.Wrap(err, "cannot parse JSON in "+string(response))
	}
	if len(eksClusters) != 1 {
		return zero, errors.New(fmt.Sprintf("Found %d clusters when searching for name %s in region %s. Found %v", len(eksClusters), clusterName, region, eksClusters))
	}
	eksClus := eksClusters[0]
	return eksClus, nil
}

func runDescribeCluster(clusterName string, region string) (response []byte, err error) {

	tempStdoutFile, oldStdout := util.ReplaceStdoutOrErr(true)
	defer util.RestoreStdoutOrError(tempStdoutFile, oldStdout, true)
	args := []string{"eksctl", "get", "cluster",
		"--name", clusterName, "--region", region, "--output", "json"}
	err = runEksctl(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for DescribeCluster")
	}
	response, err = ioutil.ReadFile(tempStdoutFile)

	if err != nil {
		return nil, errors.Wrap(err, "cannot load file with eksctl response "+tempStdoutFile)
	}
	return response, err
}

// DescribeNodeGroups ... TODO replace
func DescribeNodeGroups(clusterName, region string) ([]EKSNodeGroup, error) {
	response, err := runDescribeNodeGroups(clusterName, region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for DescribeNodeGroups")
	}

	if len(response) == 0 {
		return nil, errors.New("no output")
	}
	eksNodeGroups, err := parseNodeGroupsDescription(response)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse JSON in "+string(response))
	}
	if len(eksNodeGroups) != 1 {
		return nil, errors.New(fmt.Sprintf("Found %d clusters when searching for clusterName %s in region %s. Found %v", len(eksNodeGroups), clusterName, region, eksNodeGroups))
	}

	return eksNodeGroups, nil
}

func runDescribeNodeGroups(clusterName string, region string) (response []byte, err error) {

	tempStdoutFile, oldStdout := util.ReplaceStdoutOrErr(true)
	defer util.RestoreStdoutOrError(tempStdoutFile, oldStdout, true)
	args := []string{"eksctl", "get", "nodegroups",
		"--cluster", clusterName, "--region", region, "--output", "json"}
	err = runEksctl(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for DescribeNodeGroups")
	}
	response, err = ioutil.ReadFile(tempStdoutFile)

	if err != nil {
		return nil, errors.Wrap(err, "cannot load file with eksctl response "+tempStdoutFile)
	}
	return response, err
}

//ListClusters ... TODO replace
func ListClusters(region, labelFilter string) ([]string, error) {
	outputAsTable, err := runListClusters(region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runListClusters")
	}
	if len(outputAsTable) == 0 {
		return nil, errors.New("no output")
	}

	outputAsTableStr := string(outputAsTable)

	log.Println("Listing EKS clusters " + outputAsTableStr)
	clusterNames, err := parseClusterList(outputAsTableStr, region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse "+outputAsTableStr)
	}

	return clusterNames, nil
}

func runListClusters(region string) ([]byte, error) {

	tempStdoutFile, oldStdout := util.ReplaceStdoutOrErr(true)
	defer util.RestoreStdoutOrError(tempStdoutFile, oldStdout, true)
	args := []string{"eksctl", "get", "clusters", "--region", region, "--output", "table"}
	err := runEksctl(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot runEksCtl for ListClusters")
	}
	response, err := ioutil.ReadFile(tempStdoutFile)

	if err != nil {
		return nil, errors.Wrap(err, "cannot load file with eksctl response "+tempStdoutFile)
	}
	return response, err
}

func resetOsArgs(oldOsArgs []string) {
	os.Args = oldOsArgs
}
