package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"bldy.build/bldy/srclang/lexer"
	"bldy.build/bldy/srclang/srcutils"
)

func main() {
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		l := lexer.New(file, f)
		l.Debug()
		buf := &bytes.Buffer{}
		if err != nil {
			panic(err)
		}
		enc := srcutils.NewEncoder(buf)
		for tok := range l.Tokens {
			enc.Encode(tok)
		}
		ioutil.WriteFile(
			strings.Replace(file, filepath.Ext(file), ".dat", 1),
			buf.Bytes(),
			0755,
		)
	}
}
