package ziggy

import (
	"fmt"

	"bldy.build/build"
	"bldy.build/build/url"
)

var loaders = make(map[string]func(u *url.URL, bctx build.Context, wd string) (Package, error))

func init() {
	loaders["file"] = fileLoader
	loaders["http"] = httpLoader
	loaders["https"] = httpLoader
}

func Load(u *url.URL, bctx build.Context, wd string) (Package, error) {
	load, ok := loaders[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("%q is not a supported loader protocol", u.Scheme)
	}

	return load(u, bctx, wd)
}
