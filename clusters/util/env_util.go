package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// RootPath ...
func RootPath() string {
	wd, _ := os.Getwd()
	for wd != "" {
		files, err := ioutil.ReadDir(wd)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if "locations" == f.Name() {
				return wd
			}
		}
		if wd == "/" {
			wd = ""
		} else {
			wd = filepath.Dir(wd)
		}
	}
	return wd
}

// NoopWriter ...
type NoopWriter struct{}

// Write ...
func (m NoopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// ReplaceStdout ...
func ReplaceStdout() (tempStdoutFullPath string, oldStdout *os.File) {
	tempStdoutFullPath = filepath.Join(os.TempDir(), "stdout")
	fmt.Println("stdout is now set to", tempStdoutFullPath)
	oldStdout = os.Stdout                      // keep backup of the real stdout
	temp, err := os.Create(tempStdoutFullPath) // create temp file
	if err != nil {
		panic(err)
	}
	os.Stdout = temp
	return tempStdoutFullPath, oldStdout
}

// RestoreStdout ...
func RestoreStdout(oldStdout *os.File) {
	os.Stdout = oldStdout

}
