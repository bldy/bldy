package main

import (
	"os"
	"path/filepath"
	"strings"

	"bldy.build/bldy/script/lexer"
	"bldy.build/bldy/script/srcutils"
)

func main() {
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		l := lexer.New(file, f)
		out, err := os.OpenFile(
			strings.Replace(file, filepath.Ext(file), ".dat", 1),
			os.O_RDWR|os.O_CREATE,
			0755,
		)
		if err != nil {
			panic(err)
		}
		enc := srcutils.NewEncoder(out)
		for tok := range l.Tokens {
			enc.Encode(tok)
		}
		out.Close()
	}
}
