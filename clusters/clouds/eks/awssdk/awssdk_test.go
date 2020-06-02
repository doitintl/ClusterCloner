package awssdk

/*func TestParseNodeGroup(t *testing.T) {
	file := "test-data/awssdk-describe-nodegroup.json"

	path := util.RootPath() + "/" + file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ng, err := parseNodeGroupDescription(content)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, ng.ScalingConfig["DesiredSize"])
	assert.Equal(t, 33, ng.DiskSize)
	assert.Equal(t, "clus-nonself", ng.ClusterName)
}*/
