// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tap sort of looks like TestAnythingProtocol but it doesn't
// strictly adhere to the protocol, but it still looks nice
package tap

import (
	"fmt"
	"time"

	"bldy.build/build"
	"bldy.build/build/graph"
)

type Tap struct {
	start   time.Time
	done    chan struct{}
	workers map[string]*graph.Node
	cached  int
	fresh   int
}

func New() *Tap {
	return &Tap{
		workers: make(map[string]*graph.Node),
		done:    make(chan struct{}),
	}
}
func (t *Tap) Display(updates chan *graph.Node, workers int) {

	t.start = time.Now()

	for {
		select {
		case <-t.done:
			return
		case u := <-updates:
			x := ""
			switch u.Status {
			case build.Success:
				x += "ok"
			case build.Fail:
				x += "not ok"
			default:
				continue
			}
			if u.Cached {
				t.cached++
				fmt.Printf("%s\t%s\t(cached)\n", x, u.URL.String())
			} else {
				t.fresh++
				fmt.Printf("%s\t%s\t(%s)\n", x, u.URL.String(), time.Duration(u.End-u.Start))
			}

		}
	}
}
func (t *Tap) Cancel() {
	t.done <- struct{}{}
	fmt.Println()
	fmt.Printf("not ok\t(%s) \n", time.Since(t.start))
	fmt.Println()
	fmt.Println("======")
}

func (t *Tap) Finish() {
	fmt.Println()
	fmt.Printf("ok\t(%s) \n", time.Since(t.start))
	fmt.Printf("\t%.f%% was cached\n", (float32(t.cached)/float32(t.fresh+t.cached))*100.0)
}
