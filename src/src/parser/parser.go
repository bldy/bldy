package parser

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"bldy.build/bldy/src/ast"
	"bldy.build/bldy/src/lexer"
	"bldy.build/bldy/src/srcutils"
	"bldy.build/bldy/src/token"
)

type Parser struct {
	l     *lexer.Lexer
	state stateFn
	done  bool

	q     *Queue
	f     *ast.File
	scope ast.Scope
	ptr   ast.Node
}
type stateFn func(*Parser) stateFn

type File interface {
	io.ReadCloser
	Name() string
}

// New returns a new Parser
func New(src interface{}) *Parser {
	var r io.ReadCloser
	var name string
	switch s := src.(type) {
	case File:
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
		q: NewQueue(l.Tokens),
		l: l,
	}

	return p
}

func (p *Parser) Parse() ast.Node {
	p.run()
	return p.f
}

func (p *Parser) run() {
	for p.state = parseFile; p.state != nil && !p.done; {
		p.state = p.state(p)
	}
}

func (p *Parser) next() *token.Token { return p.q.Next() }
func (p *Parser) peek() *token.Token { return p.q.Peek() }

func caller() (call string, file string, line int) {
	var caller uintptr
	caller, file, line, _ = runtime.Caller(2)
	name := strings.Split(runtime.FuncForPC(caller).Name(), ".")
	callName := name[len(name)-1]
	return callName, file, line
}

func (p *Parser) mustGet(types ...token.Type) {
	if err := p.expect(types...); err != nil {
		call, file, line := caller()
		panic(fmt.Sprintf("%v\n%s:%d <%s>\n", err, file, line, call))
	}
}
func (p *Parser) expect(types ...token.Type) error {
	t := p.peek()
	for _, typ := range types {
		if t != nil && t.Type() == typ {
			return nil
		}
	}
	return fmt.Errorf("was expecting type <%s> got <%s> instead", types, srcutils.Encode(t))
}

func parseFile(p *Parser) stateFn {
	p.f = ast.NewFile()
	p.ptr = p.f
	p.scope = p.f
	return parseModule
}

func parseModule(p *Parser) stateFn {
	p.mustGet(token.MODULE)
	p.next()
	p.mustGet(token.IDENT)
	p.f.SetModule(string(p.next().Data()))

	return parseDeclerations
}

func parseDeclerations(p *Parser) stateFn {
	p.mustGet(token.FUNC, token.LET)
	switch p.peek().Type() {
	case token.FUNC:
		return parseFunc
	case token.LET:
		return parseLet
	default:
		return nil
	}
}

func parseFunc(p *Parser) stateFn {
	p.mustGet(token.FUNC)
	p.next()
	// Declare the function
	f := &ast.Function{}
	ptr := p.ptr.(ast.Scope)
	p.mustGet(token.IDENT)
	t := p.next()
	name := string(t.Data())

	ptr.Declare(name, f)
	f.SetName(name)
	f.Parent(p.ptr)
	p.ptr = f
	p.scope = f
	p.mustGet(token.LPAREN)
	p.next()
	return parseFuncParams
}

func parseFuncParams(p *Parser) stateFn {
	ptr := p.ptr.(*ast.Function)
	t := p.next()
	namelist := []*token.Token{}
	for ; t != nil && t.Type() != token.RPAREN; t = p.next() {
		switch p.peek().Type() {
		case token.COMMA, token.RPAREN:
			switch {
			case len(namelist) > 0:
				ptr.AddParams(t, namelist...)
				namelist = []*token.Token{}
			case len(namelist) < 1:
				namelist = append(namelist, t)
			}
			if p.peek().Type() == token.COMMA {
				p.next()
			}
		case token.IDENT:
			namelist = append(namelist, t)
		}
	}
	return nil
}

func parseLet(p *Parser) stateFn {
	return nil
}
