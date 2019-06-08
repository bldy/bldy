package ziggy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"go.starlark.net/starlark"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/url"
	"bldy.build/bldy/src/ziggy"
)

func init() {
	ziggy.Register("file", fileLoader)
}

type Package struct {
	name string
	ctx  build.Context
	u    url.URL
	p    *starlark.Program
	rt   ziggy.Runtime
}

func fileLoader(ctx build.Context, u *url.URL) (ziggy.Package, error) {
	dir, file := filepath.Split(u.Path)
	if !(file == "buildfile" || path.Ext(file) == ".bldy") {
		return nil, fmt.Errorf("load: could not find entry file %s", u)
	}
	base, err := url.Parse(dir)
	base.Scheme = "file"
	if err != nil {
		return nil, fmt.Errorf("local: %v", err)
	}
	return &Package{
		ctx: ctx.WithBase(base),
		u:   *u,
	}, nil
}

func (pkg *Package) Compile(rt ziggy.Runtime) (*starlark.Program, error) {
	u := pkg.u
	u.Scheme = ""
	u.Fragment = ""
	name := u.String()

	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("local: %v", err)
	}
	_, p, err := starlark.SourceProgram(name, f, func(s string) bool { _, ok := ziggy.LibZiggy(rt)[s]; return ok })
	if err != nil {
		return nil, fmt.Errorf("local: %v", err)
	}

	pkg.p = p

	pkg.rt = rt
	return p, nil
}

func (pkg *Package) Eval(t *starlark.Thread) (starlark.StringDict, error) {
	return pkg.p.Init(t, pkg.rt.Sys())
}

func (pkg *Package) Getbase() url.URL { return pkg.ctx.Getbase() }
