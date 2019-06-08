package build

import (
	"context"

	"bldy.build/bldy/src/build"

	"bldy.build/bldy/src/url"
)

type cliContext struct {
	context.Context
	base url.URL
}

func (c *cliContext) Getbase() (base url.URL) { return c.base }
func NewContext(ctx context.Context, u url.URL) build.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &cliContext{ctx, u}
}

func (c *cliContext) WithBase(base *url.URL) build.Context {
	return &cliContext{c, *base}
}
