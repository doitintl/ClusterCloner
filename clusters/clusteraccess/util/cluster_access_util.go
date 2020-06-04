package util

import (
	clusterutil "clustercloner/clusters/util"
	"fmt"
	"log"
)

// PrintFilteringResults ...
func PrintFilteringResults(cloud string, labelFilter map[string]string, matchedNames []string, unmatchedNames []string) {
	matchedS := "None matched"
	if len(matchedNames) > 0 {
		matchedS = fmt.Sprintf("matched: %v", matchedNames)
	}
	unmatchedS := "None were unmatched"

	if len(unmatchedNames) > 0 {
		unmatchedS = fmt.Sprintf("%v did not match", unmatchedNames)
	}
	log.Printf("In listing %s clusters, the label filter was %s; %v; %v", cloud, clusterutil.StrMapToStr(labelFilter), matchedS, unmatchedS)
}
