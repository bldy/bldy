package parser

import (
	"fmt"
	"io"
	"os"

	"bldy.build/bldy/srclang/ast"
	"bldy.build/bldy/srclang/lexer"
	"bldy.build/bldy/srclang/token"
)

type Parser struct {
	l     *lexer.Lexer
	state stateFn
	done  bool

	n ast.Node
}
type stateFn func(*Parser) stateFn

func New(src interface{}) *Parser {
	var r io.ReadCloser
	var name string
	switch s := src.(type) {
	case *os.File:
		r = s
		name = s.Name()
	case string:
		name = s
		var err error
		if r, err = os.Open(s); err != nil {
			panic(err)
		}
	}
	l := lexer.New(name, r)
	p := &Parser{
		l: l,
	}

	return p
}

func (p *Parser) run() {
	for p.state = parseFile; p.state != nil && !p.done; {
		p.state = p.state(p)
	}
}

func (p *Parser) next() *token.Token {
	if p.done {
		panic("parser is done")
	}
	t, ok := <-p.l.Tokens
	if !ok {
		p.done = true
	}
	return t
}

func parseFile(p *Parser) stateFn {
	p.n = &ast.File{}
	return parseModule
}

func parseModule(p *Parser) stateFn {
	if _, err := p.expect(token.MODULE); err != nil {
		return nil
	}
	switch ptr := p.n.(type) {
	case *ast.File:
		t, err := p.expect(token.IDENT)
		if err != nil {
			return nil
		}
		ptr.SetModule(string(t.Data()))
	}
	return nil
}

func (p *Parser) expect(types ...token.Type) (*token.Token, error) {
	t := p.next()
	for _, typ := range types {
		if t.Type() == typ {
			return t, nil
		}
	}
	return nil, fmt.Errorf("was expecting type <%s> got <%s> instead", types, t)
}

func (p *Parser) Parse() ast.Node {
	p.run()
	return p.n
}
