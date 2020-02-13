package srcutils

import (
	"bufio"
	"fmt"
	"io"

	"bldy.build/bldy/srclang/token"
)

/*
	&line,
	&col,
	&offset,
	&length,
	&xtyp,
	&data,
*/
const _fmt = "%d:%d:%d:%d:%q:%q"

func NewEncoder(w io.Writer) *Encoder { return &Encoder{w} }

type Encoder struct{ w io.Writer }

func Encode(t *token.Token) string {
	d := t.Data()
	return fmt.Sprintf(_fmt+"\n", t.Line, t.Column, t.Offset(), len(d), t.Kind(), d)
}

func (e *Encoder) Encode(t *token.Token) { fmt.Fprint(e.w, Encode(t)) }

type Decoder struct{ scanner *bufio.Scanner }

func NewDecoder(r io.Reader) *Decoder { return &Decoder{bufio.NewScanner(r)} }

func (d *Decoder) Decode() (*token.Token, error) {
	d.scanner.Scan()
	if err := d.scanner.Err(); err != nil {
		return nil, err
	}
	txt := d.scanner.Text()
	if len(txt) < 1 {
		return nil, io.EOF
	}
	var xtyp string
	var col, line, offset, length int
	var data string
	if _, err := fmt.Sscanf(string(txt), _fmt,
		&line,
		&col,
		&offset,
		&length,
		&xtyp,
		&data,
	); err != nil {
		return nil, err
	}
	// this is hackyyy but it's a test util so doesn't matter.
	typ := token.ERROR
	for i := token.Begin(); i <= token.End(); i++ {
		if i.String() == xtyp {
			typ = i
		}
	}
	return token.New(
		typ,
		[]byte(data),
		"", offset, line, col,
	), nil
}
