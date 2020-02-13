package ast

import "bldy.build/bldy/srclang/token"

type Node interface {
	Range() (*token.Position, *token.Position)
}

type File struct {
	module string
}

func (f *File) SetModule(s string) {
	f.module = s
}
func (f *File) Module() string                            { return f.module }
func (f *File) Range() (*token.Position, *token.Position) { return nil, nil }
