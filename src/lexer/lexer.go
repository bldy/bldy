package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"bldy.build/bldy/token"
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer.
type Lexer struct {
	Tokens chan token.Token // channel of scanned items
	r      io.ByteReader
	c      io.Closer
	done   bool
	name   string // the name of the input; used only for error reports
	buf    []byte
	input  string  // the line of text being scanned.
	state  stateFn // the next lexing function to enter
	line   int     // line number in input
	pos    int     // current position in the input
	start  int     // start position of this item
	width  int     // width of last rune read from input

	debug bool

	lastToken token.Token
}

func New(name string, r io.ReadCloser) *Lexer {
	l := &Lexer{
		r:      bufio.NewReader(r),
		c:      r,
		name:   name,
		line:   1,
		Tokens: make(chan token.Token),
		debug:  true,
	}
	go l.run()
	return l
}

// run runs the state machine for the Scanner.
func (l *Lexer) run() {
	for l.state = lexAny; l.state != nil; {
		l.state = l.state(l)
	}
	l.emit(token.EOF)
	l.c.Close()
	close(l.Tokens)
}

func (l *Lexer) emit(t token.Type) {
	s := l.input[l.start:l.pos]
	tok := token.Token{
		Type: t,
		Text: []byte(s),
		Position: token.Position{
			Filename: l.name,
			Line:     l.line,
			Offset:   l.start,
			Column:   l.pos,
		},
	}

	if l.debug {
		fmt.Printf("%s <%s> %q\n", tok.Position, t, tok)
	}
	if t != token.Newline {
		l.Tokens <- token.Token{
			Type: t,
			Text: []byte(s),
			Position: token.Position{
				Line:   l.line,
				Offset: l.start,
				Column: l.pos,
			},
		}
	}
	l.start = l.pos
	l.width = 0
	if t == token.Newline {
		l.line++
	}
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if !l.done && int(l.pos) == len(l.input) {
		l.loadLine()
	}
	if len(l.input) == l.start {
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// loadLine reads the next line of input and stores it in (appends it to) the input.
// (l.input may have data left over when we are called.)
// It strips carriage returns to make subsequent processing simpler.
func (l *Lexer) loadLine() {
	l.buf = l.buf[:0]
	for {
		c, err := l.r.ReadByte()
		if err != nil {
			l.done = true
			break
		}
		if c != '\r' {
			l.buf = append(l.buf, c)
		}
		if c == '\n' {
			break
		}
	}
	l.input = l.input[l.start:l.pos] + string(l.buf)
	l.pos -= l.start
	l.start = 0
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// errorf returns an error token and continues to scan.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Tokens <- token.Token{
		Type: token.Error,
		Text: []byte(fmt.Sprintf(format, args...)),
		Position: token.Position{
			Line:   l.line,
			Offset: l.start,
			Column: l.pos,
		},
	}
	return lexAny
}

func lexAny(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case isSpace(r):
			return lexSpace
		case r == eof:
			return nil
		case r == '-':
			return lexHyphen
		case isValidNumber(r):
			return lexNum
		case isString(r):
			return lexAlnum
		case r == '"', r == '\'':
			return lexQuote
		case r == '@':
			return lexID
		case r == '.':
			l.emit(token.Period)
		case r == ',':
			l.emit(token.Comma)
		case r == ':':
			l.emit(token.Colon)
		case r == ';':
			l.emit(token.Semicolon)
		case r == '{':
			l.emit(token.LeftCurly)
		case r == '[':
			l.emit(token.LeftBrac)
		case r == '(':
			l.emit(token.LeftParen)
		case r == ')':
			l.emit(token.RightParen)
		case r == '}':
			l.emit(token.RightCurly)
		case r == ']':
			l.emit(token.RightBrac)
		case r == '=':
			l.emit(token.Assign)
		case r == '$':
			return lexAnnotation
		case r == '#':
			return lexComment
		case isEndOfLine(r):
			return lexNewLine
		}
	}
	return nil
}

var builtins = map[string]token.Type{}
var constants = map[string]token.Type{
	"inf":   token.Inf,
	"-inf":  token.NegInf,
	"true":  token.True,
	"false": token.False,
	"nan":   token.NaN,
}

func init() {
	// Built in types https://capnproto.org/language.html#built-in-types
	for _, builtin := range []token.Type{
		token.Void,
		token.Bool,
		token.Int8,
		token.Int16,
		token.Int32,
		token.Int64,
		token.UInt8,
		token.UInt16,
		token.UInt32,
		token.UInt64,
		token.Float32,
		token.Float64,
		token.Text,
		token.Data,
		token.List,
	} {
		builtins[builtin.String()] = builtin
	}
	for _, builtin := range []token.Type{
		token.Struct,     // https://capnproto.org/language.html#structs
		token.Annotation, // https://capnproto.org/language.html#nesting-scope-and-aliases
		token.Import,     // https://capnproto.org/language.html#imports
		token.Using,      // https://capnproto.org/language.html#imports
		token.Enum,       // https://capnproto.org/language.html#enums
		token.Union,      // https://capnproto.org/language.html#unions
		token.Group,      // https://capnproto.org/language.html#groups
		token.Interface,  // https://capnproto.org/language.html#interfaces
		token.Const,      // https://capnproto.org/language.html#constants
	} {
		builtins[strings.ToLower(builtin.String())] = builtin
	}
}

func lexAlnum(l *Lexer) stateFn {
	for isAlphaNumeric(l.peek()) {
		l.next()
	}

	// Do we need this special case? it certainly makes
	// stuff easier but variables don't have this luxury
	// should variables start with let or var?
	tok := l.input[l.start:l.pos]
	if builtin, ok := builtins[tok]; ok {
		l.emit(builtin)
	} else if constant, ok := constants[tok]; ok {
		l.emit(constant)
	} else {
		l.emit(token.Ident)
	}
	return lexAny
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexAny
}

func lexComment(l *Lexer) stateFn {
	for !isEndOfLine(l.peek()) {
		l.next()
	}
	l.emit(token.Comment)
	return lexNewLine
}

func lexNum(l *Lexer) stateFn {
	emitee := token.Int
	for isValidNumber(l.peek()) {
		switch l.next() {
		case '.':
			emitee = token.Float
		case 'x':
			return lexHex
		}
	}
	l.emit(emitee)
	return lexAny
}

func lexNewLine(l *Lexer) stateFn {
	l.emit(token.Newline)
	l.loadLine()
	return lexAny
}

func lexID(l *Lexer) stateFn {
	// 'isValidHex' covers prettymuch both ints and hex
	for isValidHex(l.peek()) {
		l.next()
	}
	l.emit(token.ID)
	return lexAny
}

// lexHex lexes a hexadecimal
func lexHex(l *Lexer) stateFn {
	for isValidHex(l.peek()) {
		l.next()
	}
	l.emit(token.Hex)
	return lexAny
}

func lexHyphen(l *Lexer) stateFn {
	boo := l.peek()
	if boo == '>' {
		l.next()
		l.emit(token.Arrow)
	} else if isValidNumber(boo) {
		return lexNum
	} else {
		l.emit(token.Sub)
	}
	return lexAny
}
func lexQuote(l *Lexer) stateFn {
	l.backup()
	quote := l.next()
	l.ignore()
	for l.peek() != quote {
		l.next()
	}
	l.emit(token.Quote)

	if r := l.next(); r == quote {
		l.ignore()
		return lexAny
	} else {
		l.errorf("Unexpected character inside quote in position %d:%d character %q.",
			l.line,
			l.pos,
			r)
		return nil
	}
}

func lexAnnotation(l *Lexer) stateFn {
	l.emit(token.Annotation)
	return lexAlnum
}

// Helper functions

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isString(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return isString(r) || unicode.IsDigit(r)
}

func isValidNumber(r rune) bool {
	return unicode.IsDigit(r) ||
		r == '-' ||
		r == '.' ||
		r == 'x'

}
func isValidHex(r rune) bool {
	return unicode.IsDigit(r) ||
		r == 'x' ||
		r == 'X' ||
		r == 'A' ||
		r == 'a' ||
		r == 'B' ||
		r == 'b' ||
		r == 'c' ||
		r == 'C' ||
		r == 'd' ||
		r == 'D' ||
		r == 'e' ||
		r == 'E' ||
		r == 'f' ||
		r == 'F'
}