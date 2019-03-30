// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type Package struct {
	Dir  string // directory containing package sources
	Name string // package name

	// Source files
	BuildFiles []string // .bldy source files

	ctx *Context

	rules map[string]Rule
}

func (pkg *Package) Eval() error {
	thread := &starlark.Thread{
		Name:  pkg.Name,
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
	predeclared := starlark.StringDict{
		"rule": starlark.NewBuiltin("rule", pkg.newRule),
	}
	for _, file := range pkg.BuildFiles {
		data, err := pkg.ctx.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "pkg eval")
		}
		_, err = starlark.ExecFile(thread, file, data, predeclared)
		if err != nil {
			if evalErr, ok := err.(*starlark.EvalError); ok {
				log.Fatal(evalErr.Backtrace())
			}
			log.Fatal(err)
		}
	}
	return nil
}

func (pkg *Package) newRule(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	return starlark.None, nil
}

// Rule is a ziggy rule that is implemented in stardust
type Rule struct {
	name string

	ctx *Context
}

func findArg(kw starlark.Value, kwargs []starlark.Tuple) (starlark.Value, bool) {
	for i := 0; i < len(kwargs); i++ {
		if ok, err := starlark.Equal(kwargs[i].Index(0), kw); err == nil && ok {
			return kwargs[i].Index(1), true
		} else if err != nil {
			return nil, false
		}
	}
	return nil, false
}