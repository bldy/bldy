package build

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"

	"bldy.build/bldy/src/builder"
	"bldy.build/bldy/src/graph"
	"bldy.build/bldy/src/namespace/gvisor"
	"bldy.build/bldy/src/url"
	"bldy.build/bldy/src/ziggy"
	_ "bldy.build/bldy/src/ziggy/local"

	"github.com/google/subcommands"
)

type BuildCmd struct{ builder.Config }

func (*BuildCmd) Name() string     { return "build" }
func (*BuildCmd) Synopsis() string { return "builds a target" }
func (*BuildCmd) Usage() string {
	return `build src/libc#klibc
Builds a target
`
}

func (b *BuildCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&b.Fresh, "fresh", false, "use the cache or build fresh")
}

func (b *BuildCmd) Execute(c context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if len(args) != 1 {
		return subcommands.ExitUsageError
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return 3
	}
	x, err := url.Parse(wd)
	x.Scheme = "file"
	ctx := NewContext(c, *x)
	if err != nil || x == nil {
		fmt.Println(err.Error())
		return 4
	}

	z := ziggy.New(ctx, &gvisor.Runtime{})

	u, ok := args[0].(*url.URL)
	if !ok {
		return subcommands.ExitUsageError
	}

	g, err := graph.New(
		ctx,
		u,
		z.GetTask,
	)

	if err != nil {
		fmt.Println(err.Error())
		return 4

	}

	if g == nil {
		fmt.Println("nothing to build")
		return 5
	}

	workers := float64(runtime.NumCPU()) * 1.25

	bldr := builder.New(
		ctx,
		&gvisor.Runtime{},
		&b.Config,
		g,
		newNotifier(int(math.Round(workers))),
	)

	bldr.Execute(ctx, int(math.Round(workers)))

	if err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
