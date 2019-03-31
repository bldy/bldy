// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"bldy.build/build"
	"bldy.build/build/url"
	"go.starlark.net/starlarkstruct"
	"sevki.org/x/pretty"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type Package struct {
	Dir  string // directory containing package sources
	Name string // package name

	// Source files
	BuildFiles []string // .bldy source files

	ctx build.Context

	rules map[string]*Rule
}

func (z *ziggy) Load(u *url.URL) (*Package, error) {
	_, dir := filepath.Split(u.Path)
	files, err := filepath.Glob(filepath.Join(u.Path, "*.bldy"))
	if err != nil {
		return nil, errors.Wrap(err, "import")
	}
	return &Package{
		Dir:        u.Path,
		Name:       dir,
		BuildFiles: files,
		ctx:        z.ctx,
		rules:      make(map[string]*Rule),
	}, nil
}

func (pkg *Package) Eval() error {
	thread := &starlark.Thread{
		Name:  pkg.Name,
		Print: func(_ *starlark.Thread, msg string) { log.Println(msg) },
		Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			u, err := url.Parse(module)
			if err != nil {
				return nil, err
			}
			log.Println(pretty.JSON(u))

			return nil, errors.New("")
		},
	}
	predeclared := starlark.StringDict{
		"rule":   starlark.NewBuiltin("rule", pkg.newRule),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
	for _, file := range pkg.BuildFiles {
		data, err := ioutil.ReadFile(file)
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
