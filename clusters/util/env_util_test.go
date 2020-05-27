package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRootPath(t *testing.T) {

	r := RootPath()
	if !strings.Contains(r, "clustercloner") {
		t.Fatal(r)
	}
}
func TestReplaceStdout(t *testing.T) {

	fmt.Println("Before Stdout replaced")
	tempFile, old := ReplaceStdout()
	inString := "tempfile Stdout"
	fmt.Print(inString)
	RestoreStdout(old, "")
	fmt.Println("Stdout restored")
	read, err := ioutil.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}
	readS := string(read)
	assert.Equal(t, inString, readS)

}
