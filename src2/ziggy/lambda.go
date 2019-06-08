package ziggy

import (
	"fmt"
	"log"

	"bldy.build/bldy/src/build"
	"go.starlark.net/starlark"
)

type lambda struct {
	name string
	impl *starlark.Function `ziggy:"implementation"`
	rt   build.Runtime
}

func (l *lambda) String() string        { panic("not implemented") }
func (l *lambda) Type() string          { return fmt.Sprintf("<stardust.lambda %q>", l.name) }
func (l *lambda) Freeze()               {}
func (l *lambda) Truth() starlark.Bool  { return starlark.True }
func (l *lambda) Hash() (uint32, error) { panic("not implemented") }
func (l *lambda) Name() string          { return l.name }

func (l *lambda) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string
	err := starlark.UnpackArgs(l.Name(), args, kwargs, ziggyKeyName, &s)
	if err != nil {
		return nil, err
	}
	t := &Task{name: s, ar: &actionRecorder{}}
	outputs, err := starlark.Call(thread, l.impl, starlark.Tuple{newContext(s, l.rt, t)}, []starlark.Tuple{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	t.outputs = outputs
	return t, nil
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
