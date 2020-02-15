package parser

import (
	"os"
	"testing"

	"bldy.build/bldy/src/ast"
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

func TestParseFile(t *testing.T) {
	p := New("../testdata/file.src")
	file, ok := p.Parse().(*ast.File)
	if !ok {
		t.Fail()
		return
	}
	if file.Module() != "test" {
		t.Logf("%q is not %q", file.Module(), "test")
		t.Fail()
		return
	}
	if file.Decl("Bldy") == nil {
		t.Logf("%+v", file)
		t.Fail()
		return
	}
}
func TestFuncWithParams(t *testing.T) {
	p := New("../testdata/funcparams.src")
	file, ok := p.Parse().(*ast.File)
	if !ok {
		t.Fail()
		return
	}
	if file.Module() != "test" {
		t.Logf("%q is not %q", file.Module(), "test")
		t.Fail()
		return
	}
	if file.Decl("Bldy") == nil {
		t.Logf("%+v", file)
		t.Fail()
		return
	}
}
