package ziggy

import (
	"bldy.build/build/executor"
	"bldy.build/build/url"
)

type Rule struct {
	name string
	u    url.URL
}

func (r *Rule) Name() string {
	return r.name
}

func (r *Rule) Dependencies() []*url.URL {
	return nil
}

func (r *Rule) Outputs() []string {
	panic("not implemented")
}

func (r *Rule) Hash() []byte {
	return []byte{}
}

func (r *Rule) Build(*executor.Executor) error {
	return nil
}
