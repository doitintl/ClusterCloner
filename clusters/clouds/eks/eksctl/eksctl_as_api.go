package eksctl

import (
	"github.com/pkg/errors"
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

func resetOsArgs(oldOsArgs []string) {
	os.Args = oldOsArgs
}
