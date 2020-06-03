package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
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
	execReplaceStdOutOrErr(t, true)
}

func TestReplaceStderr(t *testing.T) {
	execReplaceStdOutOrErr(t, false)
}
func execReplaceStdOutOrErr(t *testing.T, isStdOut bool) {
	tempFile, old := ReplaceStdoutOrErr(isStdOut)
	inputString := "tempfile content"
	fmt.Print(inputString)
	read, err := ioutil.ReadFile(tempFile)
	assert.Nil(t, err)
	readS := string(read)

	RestoreStdoutOrError(tempFile, old, isStdOut)
	assert.Equal(t, inputString, readS)
	fmt.Println("Stdout restored")
	if _, err := os.Stat(tempFile); err == nil {
		t.Fatal(tempFile + " should not exist")
	}
}
