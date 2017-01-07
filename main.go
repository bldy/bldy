// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"bldy.build/bldy/tap"
	"bldy.build/build/builder"
	"bldy.build/build/util"

	"runtime"

	"flag"

	_ "bldy.build/build/targets/build"
	_ "bldy.build/build/targets/cc"
	_ "bldy.build/build/targets/golang"
	_ "bldy.build/build/targets/harvey"
	_ "bldy.build/build/targets/yacc"

	"sevki.org/lib/prettyprint"
)

var (
	buildVer = "version"
	usage    = `usage: build target

We require that you run this application inside a git project.
All the targets are relative to the git project. 
If you are in a subfoler we will traverse the parent folders until we hit a .git file.
`
)
var (
	disp    = flag.String("d", "tap", "only available display is tap currently")
	display Display
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		printUsage()
	}
	target := flag.Args()[0]

	switch *disp {
	case "tap":
		display = tap.New()
	}
	switch target {
	case "version":
		version()
		return
	case "force":
		os.RemoveAll(builder.BLDYCACHE)
		if len(flag.Args()) >= 2 {
			target = flag.Args()[1]
			execute(target)
		}
	case "clean":
		clean(target)
	case "query":
		target = flag.Args()[1]
		query(target)
	case "installs":
		target = flag.Args()[1]
		installs(target)
	case "hash":
		target = flag.Args()[1]
		hash(target)
	default:
		execute(target)
	}
}
func progress() {
	fmt.Println(runtime.NumCPU())
}
func printUsage() {
	fmt.Fprintf(os.Stderr, usage)
	os.Exit(1)

}
func version() {
	fmt.Printf("Build %s", buildVer)
	os.Exit(0)
}

func hash(t string) {
	c := builder.New()

	if c.ProjectPath == "" {
		fmt.Fprintf(os.Stderr, "You need to be in a git project.\n\n")
		printUsage()
	}
	fmt.Printf("%x\n", c.Add(t).HashNode())
}

func query(t string) {

	c := builder.New()

	if c.ProjectPath == "" {
		fmt.Fprintf(os.Stderr, "You need to be in a git project.\n\n")
		printUsage()
	}
	fmt.Println(prettyprint.AsJSON(c.Add(t).Target))
}
func installs(t string) {

	c := builder.New()

	if c.ProjectPath == "" {
		fmt.Fprintf(os.Stderr, "You need to be in a git project.\n\n")
		printUsage()
	}
	fmt.Println(prettyprint.AsJSON(c.Add(t).Target.Installs()))
}
func clean(t string) {
	c := builder.New()

	if c.ProjectPath == "" {
		fmt.Fprintf(os.Stderr, "You need to be in a git project.\n\n")
		printUsage()
	}
	target := c.Add(t).Target
	for file, _ := range target.Installs() {
		if err := os.Remove(filepath.Join(util.BuildOut(), file)); err != nil {
			log.Println(err)
		}
	}
}

func execute(t string) {
	c := builder.New()

	if c.ProjectPath == "" {
		fmt.Fprintf(os.Stderr, "You need to be in a git project.\n\n")
		printUsage()
	}
	c.Root = c.Add(t)
	c.Root.IsRoot = true

	if c.Root == nil {
		log.Fatal("We couldn't find the root")
	}
	cpus := int(float32(runtime.NumCPU()) * 1.25)

	// If the app hangs, there is a log.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		f, _ := os.Create("/tmp/build-crash-log.json")
		fmt.Fprintf(f, prettyprint.AsJSON(c.Root))
		os.Exit(1)
	}()

	go display.Display(c.Updates, cpus)

	go c.Execute(time.Second, cpus)
	for {
		select {
		case done := <-c.Done:
			if done.IsRoot {
				display.Finish()
				os.Exit(0)
			}
		case err := <-c.Error:
			display.Cancel()

			fmt.Println(err)
			os.Exit(1)
		case <-c.Timeout:
			log.Println("your build has timed out")
		}

	}

}
