package ziggy

import (
	"fmt"

	"bldy.build/build"
	"go.starlark.net/starlark"
)

type Context struct {
	build.Context

	name string
}

func (ctx *Context) String() string        { panic("not implemented") }
func (ctx *Context) Type() string          { return "build_context" }
func (ctx *Context) Freeze()               { panic("not implemented") }
func (ctx *Context) Truth() starlark.Bool  { panic("not implemented") }
func (ctx *Context) Hash() (uint32, error) { panic("not implemented") }
func (ctx *Context) AttrNames() []string   { panic("not implemented") }
func (ctx *Context) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(ctx.name), nil
	case "os":
		return starlark.String(ctx.Context.BLDYOS), nil
	case "arch":
		return starlark.String(ctx.Context.BLDYARCH), nil
	default:
		return nil, fmt.Errorf("%q is not a ctx attributw", name)
	}
}

func newContext(ctx build.Context, name string) starlark.Value {
	c := &Context{ctx, name}
	return c
}