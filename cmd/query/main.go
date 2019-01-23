package query

import (
	"context"
	"flag"
	"fmt"
	"io"

	"bldy.build/build/url"

	"bldy.build/build/graph"
	"github.com/google/subcommands"
	"sevki.org/x/pretty"
)

type QueryCmd struct {
	fresh bool
}

func (*QueryCmd) Name() string     { return "query" }
func (*QueryCmd) Synopsis() string { return "queries a target" }
func (*QueryCmd) Usage() string {
	return `query //<package>:<name>
{
...
}
`
}

func (q *QueryCmd) SetFlags(f *flag.FlagSet) {}

func (q *QueryCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if len(args) != 1 {
		return subcommands.ExitUsageError
	}
	u, ok := args[0].(*url.URL)
	if !ok {
		return subcommands.ExitUsageError
	}

	g, err := graph.New(u)
	if err != nil {
		fmt.Println(err.Error())
		return 4
	}
	if g == nil {
		io.WriteString(subcommands.DefaultCommander.Error, "we could not construct your graph")
	}
	fmt.Fprintln(subcommands.DefaultCommander.Output, pretty.JSON(g.Root.Target))
	return subcommands.ExitSuccess
}

type HashCmd struct{}

func (*HashCmd) Name() string     { return "hash" }
func (*HashCmd) Synopsis() string { return "prints the checksum for a target" }
func (*HashCmd) Usage() string {
	return `hash //<package>:<name>
deadbeef0012345
`
}

func (q *HashCmd) SetFlags(f *flag.FlagSet) {}

func (q *HashCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if len(args) != 1 {
		return subcommands.ExitUsageError
	}
	u, ok := args[0].(*url.URL)
	if !ok {
		return subcommands.ExitUsageError
	}

	g, err := graph.New(u)
	if err != nil {
		fmt.Println(err.Error())
		return 4
	}
	if g == nil {
		io.WriteString(subcommands.DefaultCommander.Error, "we could not construct your graph")
	}
	fmt.Fprintf(subcommands.DefaultCommander.Output, "%x\n", g.Root.HashNode())
	return subcommands.ExitSuccess
}
