package integrationtests

import (
	"clustercloner/clusters"
	"clustercloner/clusters/clouds/eks/access"
	"log"
	"testing"
)

// TestCreateAWSClusterFromFile ...
func TestCreateAWSClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/eks_clusters.json"
	CreateClusterFromFile(t, inputFile)

}

//TODO delete this test, otting the content into a big integ test
// TestCreateAWSClusterFromFile ...
func TestDescribeAWSCluster(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var err error

	ca := access.EKSClusterAccess{}
	searchTemplate := &clusters.ClusterInfo{Name: "clus-sudic", Location: "us-east-2"}
	out, err := ca.Describe(searchTemplate)
	log.Println(out)

	_ = err

}

//TODO delete this test, otting the content into a big integ test
// TestCreateAWSClusterFromFile ...
func TestDeleteCluster(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var err error

	ca := access.EKSClusterAccess{}
	searchTemplate := &clusters.ClusterInfo{Name: "clus-sudic", Location: "us-east-2"}
	err = ca.Delete(searchTemplate)
	if err != nil {
		t.Fatal(err)
	}

}
