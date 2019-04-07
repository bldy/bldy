// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"github.com/pkg/errors"

	"bldy.build/build"
	"bldy.build/build/url"
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

type ziggy struct {
	wd      string
	ctx     build.Context
	rules   map[string]build.Rule
	globals starlark.StringDict
}

func New(wd string, ctx build.Context) build.VM {
	return &ziggy{
		ctx: ctx,
		wd:  wd,
	}
}

func (z *ziggy) GetTarget(u *url.URL) (build.Rule, error) {
	pkg, err := Load(u, z.ctx, z.wd)
	if err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	if _, err := pkg.Eval(nil); err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	return pkg.GetTarget(u)
}
