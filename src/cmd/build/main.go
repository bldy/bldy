package build

import (
	"context"
	"flag"
	"net/url"

	"bldy.build/bldy/module"
	"github.com/google/subcommands"
)

type BuildCmd struct {
	Target *url.URL
	Fresh  bool
}

func (*BuildCmd) Name() string     { return "build" }
func (*BuildCmd) Synopsis() string { return "builds a target" }
func (*BuildCmd) Usage() string {
	return `bldy build src/libc
Builds a target
`
}

func (b *BuildCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&b.Fresh, "fresh", false, "use the cache or build fresh")
}

func (b *BuildCmd) Execute(c context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// get the URL that's passed as an argument
	u, ok := args[0].(*url.URL)
	if !ok {
		return subcommands.ExitUsageError
	}

	mod, err := module.New(u)
	if err != nil {
		panic(err)
	}
	_ = mod
	return subcommands.ExitSuccess
}
