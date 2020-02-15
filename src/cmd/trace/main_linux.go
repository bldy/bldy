package trace

import (
	"context"
	"flag"
	"fmt"
	"os"

	"bldy.build/bldy/trace"
	"github.com/google/subcommands"
)

type TraceCmd struct {
}

func (*TraceCmd) Name() string     { return "trace" }
func (*TraceCmd) Synopsis() string { return "traces a build" }
func (*TraceCmd) Usage() string {
	return `bldy trace -- make all
launches a sub project process and traces the syscalls it makes.
`
}
func (t *TraceCmd) SetFlags(f *flag.FlagSet) {
}

func (t *TraceCmd) Execute(c context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := []string{}
	foundDash := false
	for _, arg := range os.Args {
		if foundDash {
			args = append(args, arg)
		}
		if arg == "--" {
			foundDash = true
		}
	}

	fmt.Printf("Running %v\n", args)

	trace.Trace(args[0], args[1:]...)
	return subcommands.ExitSuccess

}
