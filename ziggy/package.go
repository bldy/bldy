// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"bldy.build/build"
	"bldy.build/build/url"

	"go.starlark.net/starlark"
)

type Package interface {
	Eval(thread *starlark.Thread) (starlark.StringDict, error)
	GetTarget(u *url.URL) (build.Rule, error)
}
