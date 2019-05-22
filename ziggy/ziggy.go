// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"io"

	"bldy.build/build"
	"bldy.build/build/url"
	"go.starlark.net/starlark"
	"golang.org/x/exp/errors"
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
	wd      string
	ctx     build.Context
	rules   map[string]build.Rule
	globals starlark.StringDict
}

func New(ctx build.Context) build.Store {
	return &ziggy{}
}

func (z *ziggy) GetTarget(u *url.URL) (build.Rule, error) {
	pkg, err := Load(u, z.ctx, z.wd)
	if err != nil {
		return nil, errors.New("get target")
	}
	if _, err := pkg.Eval(nil); err != nil {
		return nil, errors.New("get target")
	}
	return pkg.GetTarget(u)
}

func (z *ziggy) Eval(r io.Reader) error {
	return nil
}
