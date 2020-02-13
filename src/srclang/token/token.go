package token

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

//go:generate stringer -type Type
type Type int

const (
	// Special token types
	EOF Type = iota
	NEWLINE
	ERROR

	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// Keywords
	TYPE   // type MACH file
	MODULE // module
	LET    // let

	//
	IDENT  // main
	STRING // "something"
	INT
	FLOAT
	HEX
)

var tokens = [...]string{
	EOF:     "EOF",
	NEWLINE: string([]byte{'\n'}),
	ERROR:   "ERR",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	TYPE:   "type",
	MODULE: "module",
	LET:    "let",
}
var keywords map[string]Type

func caller() (call string, file string, line int) {
	var caller uintptr
	caller, file, line, _ = runtime.Caller(2)
	name := strings.Split(runtime.FuncForPC(caller).Name(), ".")
	callName := name[len(name)-1]
	return callName, file, line
}

func init() {
	keywords = make(map[string]Type)
	for i := ADD; i <= LET; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Type {
	if typ, is_keyword := keywords[ident]; is_keyword {
		return typ
	}
	return IDENT
}

func New(t Type, data []byte, file string, offset, line, column int) *Token {
	return &Token{
		t:    t,
		data: data,
		Position: Position{
			Filename: file,
			Line:     line,
			offset:   uint64(offset),
			Column:   column,
		},
	}
}

// Position describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
type Position struct {
	Filename string // filename, if any
	offset   uint64 // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

type Token struct {
	t    Type
	data []byte
	Position
}

func (t *Token) Kind() Type     { return t.t }
func (t *Token) Data() []byte   { return t.data }
func (t *Token) Len() int       { return len(t.data) }
func (t *Token) Offset() uint64 { return t.offset }

func (t *Token) Is(a *Token) error {
	if a == nil {
		panic("a cannot be nil")
	}
	if t.Kind() != a.Kind() {
		return fmt.Errorf("%s is not the same as %s", t.t, a.t)
	}
	if t.Offset() != a.Offset() {
		return fmt.Errorf("expected offset of %d got %d instead", t.Offset(), a.Offset())
	}
	if t.Line != a.Line {
		return fmt.Errorf("expected line of %d got %d instead", t.Line, a.Line)
	}
	if t.Column != a.Column {
		return fmt.Errorf("expected line of %d got %d instead", t.Column, a.Column)
	}
	if bytes.Compare(t.Data(), a.Data()) != 0 {
		return fmt.Errorf("%q is not the same as %q", t.Data(), a.Data())
	}
	return nil
}
