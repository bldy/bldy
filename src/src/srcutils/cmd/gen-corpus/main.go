package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func main() {
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		sum := sha256.Sum256(content)
		if err := ioutil.WriteFile(path.Join(os.Args[2], fmt.Sprintf("%x", sum)), content, 0644); err != nil {
			panic(err)
		}

	}
}
