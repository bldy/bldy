package lexer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"unicode"
	"unicode/utf8"

	"bldy.build/bldy/srclang/token"
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer.
type Lexer struct {
	Tokens chan *token.Token // channel of scanned items

	done bool
	name string // the name of the input; used only for error reports

	r     io.ByteReader
	c     io.Closer
	buf   []byte
	lb    *lineBuffer
	state stateFn // the next lexing function to enter

	debug     bool
	lastToken token.Token
}

// New returns a new srclang lexer
func New(name string, r io.ReadCloser) *Lexer {
	l := &Lexer{
		r:    bufio.NewReader(r),
		c:    r,
		name: name,
		lb: &lineBuffer{
			line: 0,
			col:  0,
			r:    bufio.NewReader(r),
		},
		Tokens: make(chan *token.Token),
		debug:  false,
	}
	go l.run()
	return l
}

// run runs the state machine for the Scanner.
func (l *Lexer) run() {
	for l.state = lexAny; l.state != nil && !l.done; {
		l.state = l.state(l)
	}
	l.emit(token.EOF)
	l.c.Close()
	close(l.Tokens)
}

func (l *Lexer) buffer() []byte { return l.lb.buffer() }

func (l *Lexer) emit(t token.Type) {
	s := l.lb.buffer()
	tok := token.New(t, s, l.name, int(l.lb.offset)-len(s), l.lb.line, l.lb.col-utf8.RuneCount(s)+1)
	if l.debug {
		call, file, line := caller()
		log.Println(tok)
		log.Printf("%s:%d <%s>\n", file, line, call)
	}
	l.lb.markStart()
	if t != token.NEWLINE {
		l.Tokens <- tok
	}
}

func (l *Lexer) next() rune {
	c, _, err := l.lb.ReadRune()
	if err == io.EOF {
		l.done = true
	} else if err != nil {
		l.done = true
		l.errorf("error reading line: %v", err)
		return -1
	}
	return c
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() { l.lb.markStart() }

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	if err := l.lb.UnreadRune(); err != nil {
		panic(err)
	}
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune { return l.lb.peek() }

// errorf returns an error token and continues to scan.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	s := l.lb.buffer()
	l.Tokens <- token.New(token.ERROR, []byte(fmt.Sprintf(format, args...)), l.name, int(l.lb.offset)-len(s), l.lb.line, l.lb.col)
	return lexAny
}

func lexAny(l *Lexer) stateFn {
	for !l.done {
		r := l.next()
		switch {
		case isSpace(r):
			return lexSpace
		case isEndOfLine(r):
			return lexNewLine
		case unicode.IsDigit(r):
			return lexNumber
		case isIdent(r):
			return lexIdent
		default:
			switch r {
			case '"':
				return lexString
			case '{', '}', '(', ')', '=', ':':
				return lexOperator
			}
			return nil
		}
	}
	return nil
}
func lexString(l *Lexer) stateFn {
	fr, _ := utf8.DecodeRune(l.buffer())
	l.next()
	for lr, _ := utf8.DecodeLastRune(l.buffer()); lr != fr; lr, _ = utf8.DecodeLastRune(l.buffer()) {
		l.next()
	}
	l.emit(token.STRING)
	return lexAny
}

func lexNumber(l *Lexer) stateFn {
	typ := token.INT
	for isValidNumber(l.peek()) {
		switch l.next() {
		case 'x', 'X':
			typ = token.HEX
		case '.':
			typ = token.FLOAT
		}
	}
	l.emit(typ)
	return lexAny
}

func lexType(l *Lexer) stateFn {
	l.emit(token.COLON)
	return lexAny
}

func lexOperator(l *Lexer) stateFn {
	for {
		if typ := token.Lookup(string(l.buffer())); typ != token.ERROR {
			l.emit(typ)
			return lexAny
		}

		if !isSpace(l.peek()) && !isEndOfLine(l.peek()) {
			l.next()
		} else {
			break
		}
	}

	return lexAny
}

func lexIdent(l *Lexer) stateFn {
	for isIdent(l.peek()) {
		l.next()
	}
	l.emit(token.Lookup(string(l.buffer())))
	return lexAny
}
func lexAlphaNumeric(l *Lexer) stateFn {
	for isAlphaNumeric(l.peek()) {
		l.next()
	}
	l.emit(token.Lookup(string(l.buffer())))
	return lexAny
}

func lexNewLine(l *Lexer) stateFn { l.emit(token.NEWLINE); return lexAny }

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
func isSpace(r rune) bool { return unicode.IsSpace(r) }

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isIdent reports whether r is an alphabetic, digit, or underscore.
func isIdent(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return isIdent(r) || unicode.IsDigit(r)
}

func isValidNumber(r rune) bool {
	return unicode.IsDigit(r) ||
		r == '-' ||
		r == '.' ||
		isValidHex(r)

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
