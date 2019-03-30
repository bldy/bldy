package ziggy

import (
	"path/filepath"

	"bldy.build/build/url"
	"github.com/pkg/errors"
)

type Context struct {
	BLDYARCH string // target architecture
	BLDYOS   string // target operating system

	protocol string
 
	ReadFile func(path string) ([]byte, error)
}

func (ctx *Context) Import(u *url.URL) (*Package, error) {
	_, dir := filepath.Split(u.Path)
	files, err := filepath.Glob(filepath.Join(u.Path, "*.bldy"))
	if err != nil {
		return nil, errors.Wrap(err, "import")
	}
	return &Package{
		Dir:        u.Path,
		Name:       dir,
		BuildFiles: files,
		ctx:        ctx,
	}, nil
}
