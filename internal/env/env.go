package env

import (
	"fmt"
	"io/ioutil"
	"os"
)

// FindLatestFile Returns the latest file in a directory
func FindLatestFile(path string) string {
	files, _ := ioutil.ReadDir(path)
	var newestFile string
	var newestTime int64 = 0
	for _, f := range files {
		fi, err := os.Stat(path + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = f.Name()
		}
	}

	return newestFile
}
