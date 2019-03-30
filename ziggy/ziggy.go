// Package ziggy is the is the VM for stardust scripts which are
// derivied from starlark scripts
package ziggy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"bldy.build/build/project"
	"github.com/pkg/errors"

	"bldy.build/build"
	"bldy.build/build/url"
	"go.starlark.net/starlark"
	"sevki.org/x/pretty"
)

var DefaultContext = Context{
	BLDYARCH: runtime.GOARCH,
	BLDYOS:   runtime.GOOS,
	ReadFile: ioutil.ReadFile,
}

type ziggy struct {
	wd      string
	ctx     Context
	rules   map[string]build.Rule
	globals starlark.StringDict
}

func New(wd string) build.VM {
	return &ziggy{
		ctx: DefaultContext,
		wd:  wd,
	}
}

func (z *ziggy) GetTarget(u *url.URL) (build.Rule, error) {
	if err := z.normalzieURL(u); err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	pkg, err := z.ctx.Import(u)
	if err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	if err := pkg.Eval(); err != nil {
		return nil, errors.Wrap(err, "get target")
	}
	if r, ok := z.rules[u.String()]; ok {
		return r, nil
	}

	return nil, fmt.Errorf("ziggy: %s could not be found", u.String())
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
	log.Println(pretty.JSON(u))
	return nil
}
