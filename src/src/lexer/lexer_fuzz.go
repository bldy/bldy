// +build gofuzz

package lexer

import (
	"bytes"
	"io/ioutil"

	"bldy.build/bldy/src/token"
)

// Fuzz is for fuzzing the lexer
func Fuzz(data []byte) int {
	buf := bytes.NewBuffer(data)
	l := New("fuzz", ioutil.NopCloser(buf))
	for t := range l.Tokens {
		if t.Type() == token.ERROR {
			panic(t)
		}
	}
	return 0
}
