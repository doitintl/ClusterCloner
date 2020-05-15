package util

import (
	"clustercloner/clusters"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestUnmarshall(t *testing.T) {

	fn := RootPath() + "/testdata/gke_clusters.json"
	dat, err := ioutil.ReadFile(fn)

	var cis []*clusters.ClusterInfo
	err = json.Unmarshal(dat, &cis)
	if err != nil {
		t.Error("error:", err)
	}
	assert.Equal(t, 1, len(cis))
	assert.NotEqual(t, nil, cis[0].SourceCluster)
	assert.Equal(t, clusters.GCP, cis[0].Cloud)
}
