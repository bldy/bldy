package lexer

import (
	"bufio"
	"fmt"
	"unicode/utf8"
)

type lineBuffer struct {
	r         *bufio.Reader
	s         []byte
	isPrefix  bool
	line, col int
	width     int

	i        int64 // current reading index
	prevRune int64 // index of previous rune; or < 0
	start    int64
	offset   int64
}

func (r *lineBuffer) markStart()     { r.start = r.i }
func (r *lineBuffer) buffer() []byte { return r.s[r.start:r.i] }
func (r *lineBuffer) peek() rune     { ch, _ := utf8.DecodeRune(r.s[r.i:]); return ch }
func (r *lineBuffer) first() rune    { ch, _ := utf8.DecodeRune(r.buffer()); return ch }

func (r *lineBuffer) UnreadRune() error {
	if r.i <= 0 {
		return fmt.Errorf("you can't call unread consecutively ")
	}
	r.offset -= int64(r.width)
	r.i -= int64(r.width)
	if r.col == 1 {
		r.line--
	} else {
		r.col--
	}
	r.prevRune = -1
	return nil
}

// ReadRune implements the io.RuneReader interface.
func (r *lineBuffer) ReadRune() (ch rune, size int, err error) {
RESTART:
	if r.i >= int64(len(r.s)) {
		var err error
		if r.s, r.isPrefix, err = r.r.ReadLine(); err == nil {
			r.i = 0
			r.start = 0
			r.col = 0
			r.line++
			goto RESTART
		}
		r.prevRune = -1
		r.i = 0
		r.start = 0
		return 0, 0, err
	}
	defer func() { r.col++ }()
	r.prevRune = r.i
	if c := r.s[r.i]; c < utf8.RuneSelf {
		r.width = 1
		r.i++
		r.offset++
		return rune(c), 1, nil
	}
	ch, r.width = utf8.DecodeRune(r.s[r.i:])
	r.i += int64(r.width)
	r.offset += int64(r.width)

	return
}
