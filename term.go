// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"bldy.build/build/builder"
)

type TerminalTapDisplay struct {
	start   time.Time
	done    chan struct{}
	workers map[string]*builder.Node
	cached  int
	fresh   int
}

func NewTap(workers int) *TerminalTapDisplay {
	return &TerminalTapDisplay{
		workers: make(map[string]*builder.Node, workers),
		done:    make(chan struct{}),
	}
}
func (t *TerminalTapDisplay) Display(updates chan *builder.Node, workers int) {

	t.start = time.Now()

	for {
		select {
		case <-t.done:
			return
		case u := <-updates:
			x := ""
			switch u.Status {

			case builder.Success:
				x += "ok"
			case builder.Fail:
				x += "not ok"
			default:
				continue
			}
			if u.Cached {
				t.cached++
				fmt.Printf("%s\t%s\t(cached)\n", x, u.Url.String())
			} else {
				t.fresh++
				fmt.Printf("%s\t%s\t(%s)\n", x, u.Url.String(), time.Duration(u.End-u.Start))
			}

		}
	}
}
func (t *TerminalTapDisplay) Stop() {
	log.Println("cancel")
	t.done <- struct{}{}
}
func (t *TerminalTapDisplay) Finish() {
	fmt.Println()
	fmt.Printf("ok\t(%s) \n", time.Since(t.start))
	fmt.Printf("\t%.f%% was cached\n", (float32(t.fresh+t.cached)/float32(t.cached))*100.0)
}
