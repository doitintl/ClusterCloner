package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func PrintAsJson(props interface{}) {
	jsonByteArr, err := json.MarshalIndent(props, "", " ")
	if err != nil {
		log.Print(err)
	}
	jsonStr := string(jsonByteArr)
	log.Println(jsonStr)
}

// ReadJSON reads a json file, and unmashals it.
// Very useful for template deployments.
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
