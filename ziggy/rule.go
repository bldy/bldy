package ziggy

import (
	"bldy.build/build/executor"
	"bldy.build/build/url"
)

type Rule struct {
	name string
	l    *lambda
}

func (r *Rule) Name() string {
	panic("not implemented")
}

func (r *Rule) Dependencies() []*url.URL {
	panic("not implemented")
}

func (r *Rule) Outputs() []string {
	panic("not implemented")
}

func (r *Rule) Hash() []byte {
	panic("not implemented")
}

func (r *Rule) Build(*executor.Executor) error {
	panic("not implemented")
}
