package ziggy

import (
	"fmt"

	"go.starlark.net/starlark"

	"bldy.build/bldy/src/build"
)

func newContext(name string, rt build.Runtime, t *Task) starlark.Value {
	return &Context{rt, name, t}
}

type Context struct {
	build.Runtime
	name string
	t    *Task
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
	case "ARCH":
		return starlark.String(ctx.Arch()), nil
	case "OS":
		return starlark.String(ctx.OS()), nil
	case "runtime":
		return starlark.String(fmt.Sprintf("%T", ctx.Runtime)[1:]), nil
	case "exec", "os":
		return ctx.t, nil
	default:
		return nil, fmt.Errorf("%q is not a ctx attribute", name)
	}
}
