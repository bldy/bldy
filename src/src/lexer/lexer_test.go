package lexer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bldy.build/bldy/src/srcutils"
)

func TestLex(t *testing.T) {
	files, err := filepath.Glob("../testdata/*.src")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fail()
				t.Log(err)
				return
			}
			l := New(f.Name(), f)
			//	l.Debug()
			dat, err := os.Open(strings.Replace(file, filepath.Ext(file), ".gold", 1))
			if err != nil {
				t.Fail()
				t.Log(err)
				return
			}
			dec := srcutils.NewDecoder(dat)
			for got := range l.Tokens {
				expected, err := dec.Decode()
				if err != nil {
					t.Logf("was expecting %s got error %s instead", srcutils.Encode(got), err)
					t.Fail()
					return
				}
				if err := got.Is(expected); err != nil {
					t.Logf(`
x:%s
g:%s`, srcutils.Encode(expected), srcutils.Encode(got))
					t.Log(err)
					t.Fail()
					return
				}

			}
		})
	}

}
