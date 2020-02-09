package lexer

import (
	"bldy.build/bldy/script/srcutils"
	"os"
	"testing"
)

func TestLex(t *testing.T) {
	f, err := os.Open("../testdata/lexer.bldy")
	if err != nil {
		t.Fail()
		t.Log(err)
		return
	}
	l := New("testlexer", f)
	enc := srcutils.NewEncoder(os.Stdout)
	for tok := range l.Tokens {
		enc.Encode(tok)
	}
}
