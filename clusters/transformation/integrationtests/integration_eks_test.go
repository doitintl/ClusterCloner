package integrationtests

import (
	"testing"
)

// TestCreateAWSClusterFromFile ...
func TestCreateAWSClusterFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var inputFile = "test-data/eks_clusters.json"
	execTestClusterFromFile(t, inputFile)

}

/*//TODO delete (temp)
func TestDescCluster(t *testing.T) {

	ca := clusteraccess.GetClusterAccess(clusters.AWS)
	searchTemplate := &clusters.ClusterInfo{Name: "clus-bumping", Location: "us-east-2"}
	cis, err := ca.Describe(searchTemplate)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(cis)
}
*/

/*TODO delete (temp)
func TestDelCluster(t *testing.T) {

	ca := clusteraccess.GetClusterAccess(clusters.AWS)
	searchTemplate := &clusters.ClusterInfo{Name: "clus-outlung", Location: "us-east-2"}
	err := ca.Delete(searchTemplate)
	if err != nil {
		t.Fatal(err)
	}
}*/
