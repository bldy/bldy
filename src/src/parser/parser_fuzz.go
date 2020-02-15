// +build gofuzz

package parser

import (
	"bytes"
)

type fuzzFile struct {
	*bytes.Buffer
}

func (*fuzzFile) Name() string { return "fuzz.src" }

// Fuzz is for fuzzing the lexer
func Fuzz(data []byte) int {
	buf := bytes.NewBuffer(data)
	p := New(&fuzzFile{buf})
	f := p.Parse()
	if f == nil {
		panic("")
	}
	return 0
}
