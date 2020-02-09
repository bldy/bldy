package srcutils

import (
	"bufio"
	"fmt"
	"io"

	"bldy.build/bldy/script/token"
)

const _fmt = "%d:%d:%d:%s\n"

// offset
// kind
// len
// data
func NewEncoder(w io.Writer) *Encoder { return &Encoder{w} }

type Encoder struct{ w io.Writer }

func (e *Encoder) Encode(t *token.Token) {
	d := t.Data()
	fmt.Fprintf(e.w, _fmt, t.Offset(), t.Kind(), len(d), d)
}

type Decoder struct{ scanner *bufio.Scanner }

func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	return &Decoder{scanner}
}

func (d *Decoder) Decode() (*token.Token, error) {
	d.scanner.Scan()
	var typ token.Type
	var offset, l int
	var data string
	if _, err := fmt.Sscanf(string(d.scanner.Text()), _fmt,
		&offset,
		&typ,
		&l,
		&data,
	); err != nil {
		return nil, err
	}
	return token.New(
		typ,
		[]byte(data),
		"", offset, 0, 0,
	), nil
}

type Checker struct{ r io.Reader }

func (c *Checker) Check(t interface{}) {}
func NewChecker(file string) *Checker {
	return nil
}
