package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// MarshallToJSONString ...
func MarshallToJSONString(props interface{}) string {
	jsonByteArr, err := json.MarshalIndent(props, "", "  ")
	if err != nil {
		log.Print(err)
		return "<ERROR>"
	}
	jsonStr := string(jsonByteArr)
	return jsonStr
}

// ReadJSON ...
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read template file: %v\n", err)
	}
	contents := make(map[string]interface{})
	if err := json.Unmarshal(data, &contents); err != nil {
		return nil, err
	}
	return &contents, nil
}
