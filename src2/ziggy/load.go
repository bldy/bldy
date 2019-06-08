package ziggy

import (
	"fmt"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/url"
)

type PackageLoader func(ctx build.Context, u *url.URL) (Package, error)

var loaders = make(map[string]PackageLoader)

func Register(scheme string, loader PackageLoader) {
	loaders[scheme] = loader
}

func loadPackage(ctx build.Context, u *url.URL) (Package, error) {
	if !u.IsAbs() {
		base := ctx.Getbase()
		var err error
		u, err = base.Append(u)
		if err != nil {
			return nil, fmt.Errorf("load: %v", err)
		}
	}
	load, ok := loaders[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("%q does not have a scheme", u)
	}
	return load(ctx, u)
}
