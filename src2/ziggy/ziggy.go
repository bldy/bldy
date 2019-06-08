// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"fmt"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/url"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ziggyKeyImpl    = "implementation"
	ziggyKeyAttrs   = "attrs?"
	ziggyKeyDeps    = "deps?"
	ziggyKeyOutputs = "outputs?"
	ziggyKeyName    = "name"
	ziggyKeyCtx     = "ctx"
)

type ziggy struct {
	ctx     build.Context
	rt      Runtime
	tasks   map[string]build.Task
	globals starlark.StringDict
}

type Runtime interface {
	build.Runtime

	Sys() starlark.StringDict
}

func New(ctx build.Context, rt Runtime) build.Store {
	return &ziggy{
		ctx: ctx,
		rt:  rt,
	}
}

func evalpkg(ctx build.Context, rt Runtime, u *url.URL) (starlark.StringDict, error) {
	pkg, err := loadPackage(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("ziggy: get task: %v", err)
	}

	p, err := pkg.Compile(rt)
	if err != nil {
		return nil, fmt.Errorf("compile: %v", err)
	}
	base := pkg.Getbase()

	thread := &starlark.Thread{
		Name: u.String(),
		Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			modurl, err := url.Parse(module)
			if err != nil {
				return nil, fmt.Errorf("loading %q failed: %v", module, err)
			}
			return evalpkg(ctx.WithBase(&base), rt, modurl)
		},
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
	return p.Init(thread, LibZiggy(rt))
}

func (z *ziggy) GetTask(u *url.URL) (build.Task, error) {
	g, err := evalpkg(z.ctx, z.rt, u)
	_ = g
	if err != nil {
		return nil, fmt.Errorf("ziggy getting %q failed: %v", u.String(), err)
	}
	if t, ok := g[u.Fragment].(*Task); ok {
		t.name = u.String()
		return t, nil
	}
	return nil, fmt.Errorf("ziggy: getting %q failed", u.String())
}

func LibZiggy(rt build.Runtime) starlark.StringDict {
	stdlib := make(starlark.StringDict)
	stdlib["struct"] = starlark.NewBuiltin("struct", starlarkstruct.Make)
	stdlib["module"] = starlark.NewBuiltin("module", starlarkstruct.MakeModule)

	stdlib["task"] = starlark.NewBuiltin("task", taskfactory(rt))
	stdlib["exec"] = starlark.None

	return stdlib
}

func taskfactory(rt build.Runtime) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var f *starlark.Function
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, ziggyKeyImpl, &f); err != nil {
			return nil, err
		}
		return &lambda{impl: f, rt: rt}, nil
	}
}
