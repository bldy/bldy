// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "bldy.build/build/builder"

type Display interface {
	Display(chan *builder.Node, int)
	Cancel()
	Finish()
}
