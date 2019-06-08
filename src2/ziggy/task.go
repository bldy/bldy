package ziggy

import (
	"bldy.build/bldy/src/executor"
	"bldy.build/bldy/src/racy"
	"bldy.build/bldy/src/url"
	"go.starlark.net/starlark"
)

type Task struct {
	name    string `ziggy:"name"`
	u       url.URL
	ar      *actionRecorder
	outputs starlark.Value
}

func (t *Task) Name() string             { return t.name }
func (t *Task) Dependencies() []*url.URL { return nil }
func (t *Task) Outputs() []string        { return []string{} }

func (t *Task) Sum() []byte {
	h := racy.New()
	h.HashStrings(t.name)
	return racy.XOR(h.Sum(nil), t.ar.sum)
}

func (t *Task) Run(e *executor.Executor) error {
	for _, action := range t.ar.calls {
		if err := action.Do(e); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) String() string        { panic("not implemented") }
func (t *Task) Type() string          { return "*ziggy.Task" }
func (t *Task) Freeze()               { panic("not implemented") }
func (t *Task) Truth() starlark.Bool  { panic("not implemented") }
func (t *Task) Hash() (uint32, error) { panic("not implemented") }
func (t *Task) AttrNames() []string   { panic("not implemented") }
func (t *Task) Attr(name string) (starlark.Value, error) {
	return starlark.NewBuiltin(name, t.ar.newAction(name)), nil
}
func (t *Task) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	panic("not implemented")
}
