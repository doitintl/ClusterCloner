package eksctl

import (
	"clustercloner/clusters/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// TestParseClusterList ...
func TestParseClusterList(t *testing.T) {
	s := "NAME\t\tREGION\nclus-sudic\tus-east-2\nclus-2\tus-east-2\n"
	parsed, err := parseClusterList(s, "us-east-2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(parsed), 2)
	assert.Equal(t, parsed[0], "clus-sudic")
}

func TestParseClusterDescription(t *testing.T) {
	file := "test-data/eks-describecluster.json"
	path := util.RootPath() + "/" + file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	eksCluster, err := parseClusterDescription(content)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(eksCluster), 1)
	assert.Equal(t, eksCluster[0].Name, "clus-sudic")
	assert.Equal(t, eksCluster[0].Version, "1.15")
}
