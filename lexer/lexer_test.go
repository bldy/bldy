// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer // import "sevki.org/build/lexer"

import (
	"fmt"
	"os"
	"testing"

	"sevki.org/build/token"
)

func TestMap(t *testing.T) {
	ks, _ := os.Open("map.BUILD")
	l := New("sq", ks)
	for {
		tok := <-l.Tokens
		if tok.Type != token.Newline {
			fmt.Printf("%s => %s\n", tok.Type, tok.Text)
		}
		if tok.Type == token.EOF || tok.Type == token.Error {
			break
		}
	}

}
