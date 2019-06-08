package token

import "fmt"

//go:generate stringer -type Type

type Type int

const (
	EOF Type = iota
	Newline
	Error

	ID
	String
	Space
	Int
	Float
	Hex
	Semicolon
	Comment
	LeftCurly
	RightCurly
	Equal
	Colon
	Comma
	LeftParen
	RightParen
	LeftBrac
	RightBrac
	Sub
	Add
	Quote
	Assign
	Arrow
	Period
	Ident

	Void
	Bool
	Int8
	Int16
	Int32
	Int64
	UInt8
	UInt16
	UInt32
	UInt64
	Float32
	Float64
	Text
	Data
	List

	Struct
	Field
	Using
	Import
	Union
	Group
	Enum
	Annotation
	Interface
	Const

	Inf
	NegInf
	True
	False
	NaN
)

type Token struct {
	Type Type
	Text []byte
	Position
}

func (t Token) String() string {
	if t.Type == Newline {
		return "CR"
	}
	return string(t.Text)
}

// Position describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
//
type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

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