package ast

import (
	"log"

	"bldy.build/bldy/src/srcutils"
	"bldy.build/bldy/src/token"
)

type Node interface {
	Range() (*token.Position, *token.Position)
	Parent(Node)
}

type Scope interface {
	Node
	Declare(s string, d Decleration)
}

type File struct {
	module       string
	declerations []Decleration
}

func NewFile() *File {
	return &File{
		module: "",
	}
}
func (f *File) Parent(Node) {}
func (f *File) SetModule(s string) {
	f.module = s
}
func (f *File) Module() string                            { return f.module }
func (f *File) Range() (*token.Position, *token.Position) { return nil, nil }
func (f *File) Declare(s string, d Decleration)           { f.declerations = append(f.declerations, d) }
func (f *File) Decl(s string) Decleration {
	for _, d := range f.declerations {
		if d.Name() == s {
			return d
		}
	}
	return nil
}

type Function struct {
	name   string
	stmnts []Node
	parent Node
}

func (f *Function) SetName(s string)                          { f.name = s }
func (f *Function) Parent(n Node)                             { f.parent = n }
func (f *Function) Declare(s string, d Decleration)           {}
func (f *Function) Range() (*token.Position, *token.Position) { return nil, nil }
func (f *Function) Name() string                              { return f.name }
func (f *Function) AddParams(typ *token.Token, names ...*token.Token) {
	log.Println(srcutils.Encode(typ))
}

type Type struct{ begin, end token.Token }

type Decleration interface {
	Node
	Name() string
}
