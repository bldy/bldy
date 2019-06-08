// Copyright 2018 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"flag"
	"os"

	"bldy.build/bldy/cmd/build"
	"bldy.build/bldy/cmd/trace"
	"bldy.build/bldy/fileutils"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&build.BuildCmd{}, "")
	subcommands.Register(&trace.TraceCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	if u, err := fileutils.ResolveFromWD(flag.Arg(1)); err == nil {
		os.Exit(int(subcommands.Execute(ctx, u)))
	} else {
		os.Exit(int(subcommands.Execute(ctx)))
	}
}
