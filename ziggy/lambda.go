package ziggy

import (
	"fmt"

	"bldy.build/build"
	"go.starlark.net/starlark"
)

const (
	ziggyKeyImpl    = "implementation"
	ziggyKeyAttrs   = "attrs?"
	ziggyKeyDeps    = "deps?"
	ziggyKeyOutputs = "outputs?"
	ziggyKeyName    = "name"
	ziggyKeyCtx     = "ctx"
)

// Rule is a ziggy rule that is implemented in stardust
type lambda struct {
	name string
	impl *starlark.Function
	ctx  build.Context

	register func(name string) error
}

func (l *lambda) String() string        { panic("not implemented") }
func (l *lambda) Type() string          { return fmt.Sprintf("<stardust.lambda %q>", l.name) }
func (l *lambda) Freeze()               {}
func (l *lambda) Truth() starlark.Bool  { return starlark.True }
func (l *lambda) Hash() (uint32, error) { panic("not implemented") }
func (l *lambda) Name() string          { return l.name }

func (l *lambda) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String

	if err := starlark.UnpackArgs(l.Name(), args, kwargs, ziggyKeyName, &name); err != nil {
		return nil, err
	}

	if val, err := starlark.Call(thread, l.impl, []starlark.Value{newContext(l.ctx, string(name))}, nil); err != nil {
		return val, err
	}

	l.register(string(name))
	return starlark.None, nil
}

func (pkg *Package) newRule(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var impl *starlark.Function
	attrs := new(starlark.Dict)
	outputs := new(starlark.Dict)
	var name starlark.String

	if err := starlark.UnpackArgs("ziggy.newRule", args, kwargs, ziggyKeyImpl, &impl, ziggyKeyAttrs, &attrs, ziggyKeyOutputs, &outputs); err != nil {
		return nil, err
	}
	l := &lambda{
		impl: impl,
		ctx:  pkg.ctx,
	}
	l.register = func(s string) error {
		pkg.rules[s] = &Rule{
			l:    l,
			name: string(name),
		}
		return nil
	}
	return l, nil
}

func findArg(kw starlark.Value, kwargs []starlark.Tuple) starlark.Value {
	for i := 0; i < len(kwargs); i++ {
		if ok, err := starlark.Equal(kwargs[i].Index(0), kw); err == nil && ok {
			return kwargs[i].Index(1)
		} else if err != nil {
			return nil
		}
	}
	return nil
}