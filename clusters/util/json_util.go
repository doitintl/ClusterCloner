package util

import (
	"encoding/json"
	"log"
)

// ToJSON ...
func ToJSON(props interface{}) string {
	jsonByteArr, err := json.MarshalIndent(props, "", "  ")
	if err != nil {
		log.Println(err)
		return "<ERROR>"
	}
	jsonStr := string(jsonByteArr)
	return jsonStr
}
