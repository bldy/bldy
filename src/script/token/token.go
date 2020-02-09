package token

import (
	"fmt"
)

//go:generate stringer -type Type

type Type int

const (
	EOF Type = iota
	Newline
	Error
	Module
	Ident
	String
)

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

func (t *Token) Kind() int      { return int(t.t) }
func (t *Token) Data() []byte   { return t.data }
func (t *Token) Len() int       { return len(t.data) }
func (t *Token) Offset() uint64 { return t.offset }
