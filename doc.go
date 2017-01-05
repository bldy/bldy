// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
bldy is a fast, concurrent efficient build tool.

It uses BUILD files to create a build graph, which is then executed
concurrently and independent of each other when it is able to do so.

If a cached version of the output is available it uses that instead of
doing the work again.

bldy cli interface is structured like so

	bldy [force] [-p=tap] url

a bldy url is consists of two parts, package and target

	//sys/src/cmd:date

the

	//sys/src/cmd

is the package relative to the root of the folder and the : seperated bit

	date

is the target in that package that you want to compile.

To force a full rebuild, just run bldy with force:

	bldy force //:harvey



A build file consists of couple of things build targets,

	cc_binary(
		name="hello",
		src=[
			"hello.c",
		],
	)

or built-in functions

	load("//sys/src/harvey.BUILD", "CFLAGS", "harvey_binary")

for convinience bldy has a bunch of other things you can use like macros
and for loops.

Macros look like regular build targets but they are assigned to names

	harvey_binary = cc_binary(
		copts=LIB_COMPILER_FLAGS,
		includes=[
			"//sys/include",
			"//amd64/include",
		],
		deps=CMD_DEPS,
		strip=true,
		linkopts=CMD_LINK_OPTS
	)

when you combine macros with for loops you can do even more

	[harvey_binary(
		name=c[:-2],
		srcs=[c],
	) for c in CMD_SRCS]

bldy where possible tries to be compatible with both bazel and buck.
If it's add odds with one it's probably going to be because of some incompatible
behaviour between buck and bazel and bazel is the tie breaker.

*/
package main // import "bldy.build/bldy"
