package eksctl

import (
	"os"
	"strconv"
)

// CreateCluster ...
func CreateCluster(clusterName, region, k8sVersion, tagsCsv string) error {
	err := createClusterNoNodeGroup(clusterName, region, k8sVersion, tagsCsv)
	return err
}

func createClusterNoNodeGroup(clusterName, region, k8sVersion, tagsCsv string) error {
	oldArgs := os.Args[:]
	defer resetOsArgs(oldArgs)
	os.Args = []string{"eksctl", "create", "cluster", "--without-nodegroup", "--name", clusterName, "--region", region, "--tags", tagsCsv, "--version", k8sVersion}
	return runEksctl()
}

// AddLogging ...
func AddLogging(clusterName, region, k8sVersion, tagsCsv string) error {
	oldArgs := os.Args[:]
	defer resetOsArgs(oldArgs)
	os.Args = []string{"eksctl", "utils", "update-cluster-logging", "--cluster", clusterName, "--region", region}
	return runEksctl()
}

// CreateNodeGroup ...
func CreateNodeGroup(clusterName, nodeGroupName, region, k8sVersion, nodeInstanceType, tagsCsv string, nodeCount, diskSizeGB int) error {
	oldArgs := os.Args[:]
	defer resetOsArgs(oldArgs)

	os.Args = []string{"eksctl", "create", "nodegroup", "--managed",
		"--cluster", clusterName, "--name", nodeGroupName, "--region", region, "--version", k8sVersion,
		"--node-type", nodeInstanceType, "--tags", tagsCsv,
		"--nodes", strconv.Itoa(nodeCount), "--node-volume-size", strconv.Itoa(diskSizeGB),
	}
	return runEksctl()
}

func resetOsArgs(oldOsArgs []string) {
	os.Args = oldOsArgs
}
