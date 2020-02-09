package build

import (
	"context"
	"flag"
<<<<<<< HEAD
	"net/url"

	"bldy.build/bldy/module"
	"github.com/google/subcommands"
)

type BuildCmd struct {
	Target *url.URL
	Fresh  bool
}
=======
	"fmt"
	"math"
	"os"
	"runtime"

	"bldy.build/bldy/src/builder"
	"bldy.build/bldy/src/graph"
<<<<<<< HEAD:src2/cmd/build/main.go
	"bldy.build/bldy/src/namespace/gvisor"
	"bldy.build/bldy/src/url"
	"bldy.build/bldy/src/ziggy"
	_ "bldy.build/bldy/src/ziggy/local"

=======
	"bldy.build/bldy/src/url"
>>>>>>> 97e98155e24d9c7de236ebaf33a5557c36660e2d:src/cmd/build/main.go
	"github.com/google/subcommands"
)

type BuildCmd struct{ builder.Config }
>>>>>>> 97e98155e24d9c7de236ebaf33a5557c36660e2d

func (*BuildCmd) Name() string     { return "build" }
func (*BuildCmd) Synopsis() string { return "builds a target" }
func (*BuildCmd) Usage() string {
<<<<<<< HEAD
	return `bldy build src/libc
=======
	return `build src/libc#klibc
>>>>>>> 97e98155e24d9c7de236ebaf33a5557c36660e2d
Builds a target
`
}

func (b *BuildCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&b.Fresh, "fresh", false, "use the cache or build fresh")
}

func (b *BuildCmd) Execute(c context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
<<<<<<< HEAD
	// get the URL that's passed as an argument
=======
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

>>>>>>> 97e98155e24d9c7de236ebaf33a5557c36660e2d
	u, ok := args[0].(*url.URL)
	if !ok {
		return subcommands.ExitUsageError
	}

<<<<<<< HEAD
	mod, err := module.New(u)
	if err != nil {
		panic(err)
	}
	_ = mod
=======
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
>>>>>>> 97e98155e24d9c7de236ebaf33a5557c36660e2d
	return subcommands.ExitSuccess
}
