package ziggy

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"bldy.build/build"
	"bldy.build/build/url"
	"github.com/pkg/errors"
)

type localPackage struct {
	Dir  string // directory containing package sources
	Name string // package name

	// Source files
	BuildFiles []string // .bldy source files

	ctx build.Context

	rules map[string]*Rule
}

func (pkg *localPackage) Eval(thread *starlark.Thread) (starlark.StringDict, error) {
	if thread == nil {
		thread = &starlark.Thread{
			Name:  pkg.Name,
			Print: func(_ *starlark.Thread, msg string) { log.Println(msg) },
			Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
				u, err := url.Parse(module)
				if err != nil {
					return nil, err
				}
				p, err := Load(u, pkg.ctx)
				if err != nil {
					return nil, err
				}
				return p.Eval(thread)
			},
		}
	}
	predeclared := starlark.StringDict{
		"rule":   starlark.NewBuiltin("rule", pkg.newRule),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
	global := make(starlark.StringDict)
	for _, file := range pkg.BuildFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "pkg eval")
		}
		_, err = starlark.ExecFile(thread, file, data, predeclared)
		if err != nil {
			if evalErr, ok := err.(*starlark.EvalError); ok {
				log.Fatal(evalErr.Backtrace())
			}
			log.Fatal(err)
		}
	}
	return global, nil
}

func (pkg *localPackage) newRule(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var impl *starlark.Function
	attrs := new(starlark.Dict)
	outputs := new(starlark.Dict)
	var name starlark.String

	if err := starlark.UnpackArgs("ziggy.newRule", args, kwargs, ziggyKeyImpl, &impl, ziggyKeyAttrs, &attrs, ziggyKeyOutputs, &outputs); err != nil {
		return nil, err
	}
	l := &lambda{
		impl: impl,
		ctx:  pkg.ctx,
	}
	l.register = func(s string) error {
		pkg.rules[s] = &Rule{
			l:    l,
			name: string(name),
		}
		return nil
	}
	return l, nil
}

func fileLoader(u *url.URL, bctx build.Context) (Package, error) {
	_, dir := filepath.Split(u.Path)
	files, err := filepath.Glob(filepath.Join(u.Path, "*.bldy"))
	if err != nil {
		return nil, errors.Wrap(err, "file loader")
	}
	return &localPackage{
		Dir:        u.Path,
		Name:       dir,
		BuildFiles: files,
		ctx:        bctx,
		rules:      make(map[string]*Rule),
	}, nil
}
func (pkg *localPackage) GetTarget(u *url.URL) (build.Rule, error) {
	if rule, ok := pkg.rules[u.Fragment]; ok {
		return rule, nil
	}
	return nil, fmt.Errorf("couldn't find rule %q", u)
}