package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"bldy.build/bldy/script/token"
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer.
type Lexer struct {
	Tokens chan *token.Token // channel of scanned items
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

	debug     bool
	lastToken token.Token
}

func New(name string, r io.ReadCloser) *Lexer {
	l := &Lexer{
		r:      bufio.NewReader(r),
		c:      r,
		name:   name,
		line:   1,
		Tokens: make(chan *token.Token),
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
	tok := token.New(t, []byte(s), l.name, l.start, l.line, l.pos)

	if l.debug {
	//	fmt.Println(tok)
	}
	if t != token.Newline {
		l.Tokens <- tok
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
	l.Tokens <- token.New(token.Error, []byte(fmt.Sprintf(format, args...)), l.name, l.start, l.line, l.pos)

	return lexAny
}

func lexAny(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case isString(r):
			return lexAlphaNumeric
		case isSpace(r):
			return lexSpace
		default:
			return nil
		}
	}
}

func lexAlphaNumeric(l *Lexer) stateFn {
	for isString(l.peek()) {
		l.next()
	}

	switch l.input[l.start:l.pos] {
	case "module":
		l.emit(token.Module)
		break
	default:
		l.emit(token.String)
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
