// Copyright 2018 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main // import "bldy.build/build/cmd/bldy"
import (
	"context"
	"flag"
	"net/url"
	"os"

	"bldy.build/build/cmd/build"
	"bldy.build/build/cmd/query"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&build.BuildCmd{}, "")
	subcommands.Register(&query.QueryCmd{}, "")
	subcommands.Register(&query.HashCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	if u, err := url.Parse(flag.Arg(1)); err == nil {
		os.Exit(int(subcommands.Execute(ctx, u)))
	} else {
		os.Exit(int(subcommands.Execute(ctx)))
	}
}
