// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"bldy.build/bldy/src/url"
	"go.starlark.net/starlark"
)

type Package interface {
	Compile(rt Runtime) (*starlark.Program, error)
	Eval(*starlark.Thread) (starlark.StringDict, error)
	Getbase() url.URL
}
