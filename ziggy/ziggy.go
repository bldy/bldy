// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"os"
	"path"

	"bldy.build/build/project"
	"github.com/pkg/errors"

	"bldy.build/build"
	"bldy.build/build/url"
	"go.starlark.net/starlark"
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
	if err := z.normalzieURL(u); err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	pkg, err := Load(u, z.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	if _, err := pkg.Eval(nil); err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	return pkg.GetTarget(u)
}

func (z *ziggy) normalzieURL(u *url.URL) error {
	if u.Host == project.RootKey {
		rootdir, err := project.Search(z.wd, func(s string) (os.FileInfo, error) {
			for _, ext := range []string{".git"} {
				if fi, err := os.Stat(path.Join(s, ext)); err != os.ErrNotExist {
					return fi, nil
				}
			}
			return nil, os.ErrNotExist
		})
		if err != nil {
			return err
		}
		u.Host = ""
		u.Path = path.Join(rootdir, u.Path)
		return nil
	}
	return nil
}
