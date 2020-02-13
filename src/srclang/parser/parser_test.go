package parser

import (
	"os"
	"testing"

	"bldy.build/bldy/srclang/ast"
)

func TestNewWithReader(t *testing.T) {
	f, err := os.Open("../testdata/assignment.src")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	_ = New(f)
}

func TestNewWithFileName(t *testing.T) {
	_ = New("../testdata/assignment.src")
}

func TestNewParseFile(t *testing.T) {
	p := New("../testdata/file.src")
	file, ok := p.Parse().(*ast.File)
	if !ok {
		t.Fail()
	}
	if file.Module() != "test" {
		t.Fail()
	}
}
