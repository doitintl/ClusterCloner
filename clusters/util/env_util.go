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

// ReplaceStdout ...
func ReplaceStdout() (tempStdoutFullPath string, oldStdout *os.File) {
	tempStdoutFullPath = filepath.Join(os.TempDir(), "stdout")
	oldStdout = os.Stdout // keep backup of the real stdout
	temp, err := os.Create(tempStdoutFullPath)
	if err != nil {
		panic(err)
	}
	os.Stdout = temp
	return tempStdoutFullPath, oldStdout
}

// RestoreStdout ...
func RestoreStdout(oldStdout *os.File, tempFile string) {
	os.Stdout = oldStdout
	err := os.Remove(tempFile)
	if err != nil {
		log.Println("Error removing " + tempFile + " " + err.Error())
	}

}
