package ziggy

import (
	"bldy.build/bldy/src/executor"
	"bldy.build/bldy/src/racy"
	"bldy.build/bldy/src/ziggy/ziggyutils"

	"go.starlark.net/starlark"
)

type actionRecorder struct {
	calls []executor.Action
	sum   []byte
}

func (ar *actionRecorder) Record(a executor.Action) {
	ar.sum = racy.XOR(ar.sum, a.Sum())
	ar.calls = append(ar.calls, a)
}

func (ar *actionRecorder) newAction(name string) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var i executor.Action
		switch name {
		case "Run":
			i = &run{}
		}
		if err := ziggyutils.UnpackStruct(i, kwargs); err != nil {
			return starlark.None, err
		}
		ar.Record(i)
		return starlark.None, nil
	}
}
