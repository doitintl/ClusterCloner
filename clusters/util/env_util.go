package util

import (
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

// ReplaceStdoutOrErr ...
func ReplaceStdoutOrErr(isStdOut bool) (tempFilePath string, oldStdoutOrErr *os.File) {
	var tmpFileName string
	if isStdOut {
		oldStdoutOrErr = os.Stdout
		tmpFileName = "stdout"
	} else {
		oldStdoutOrErr = os.Stderr
		tmpFileName = "stderr"
	}
	tempFilePath = filepath.Join(os.TempDir(), tmpFileName)

	temp, err := os.Create(tempFilePath)
	if err != nil {
		panic(err)
	}
	os.Stdout = temp
	return tempFilePath, oldStdoutOrErr
}

// RestoreStdoutOrError ...
func RestoreStdoutOrError(tempFile string, oldStdoutOrErr *os.File, isStdout bool) {
	if isStdout {
		os.Stdout = oldStdoutOrErr
	} else {
		os.Stderr = oldStdoutOrErr
	}
	err := os.Remove(tempFile)
	if err != nil {
		log.Println("Error removing " + tempFile + " " + err.Error())
	}

}
