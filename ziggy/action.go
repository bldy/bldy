package ziggy

import (
	"bldy.build/build/executor"
)

type actionRecorder struct {
	calls []executor.Action
}

func (ar *actionRecorder) Record(a executor.Action) {
	ar.calls = append(ar.calls, a)
}
