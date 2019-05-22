// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package build defines build target and build context structures
package build

import (
	"io"
	"runtime"

	"bldy.build/build/executor"
	"bldy.build/build/url"
)

type Context struct {
	BLDYARCH string // target architecture
	BLDYOS   string // target operating system

}

var DefaultContext = Context{
	BLDYARCH: runtime.GOARCH,
	BLDYOS:   runtime.GOOS,
}

//go:generate stringer -type=Status
// Status represents a nodes status.
type Status int

const (
	// Success is success
	Success Status = iota
	// Fail is a failed job
	Fail
	// Pending is a pending job
	Pending
	// Started is a started job
	Started
	// Fatal is a fatal crash
	Fatal
	// Warning is a job that has warnings
	Warning
	// Building is a job that's being built
	Building
)

// Rule defines the interface that rules must implement for becoming build targets.
type Rule interface {
	Name() string
	Dependencies() []*url.URL
	Outputs() []string
	Hash() []byte
	Build(*executor.Executor) error
}

// VM seperate the parsing and evauluating targets logic from rest of bldy
// so we can implement and use new grammars like jsonnet or go it self.
type Store interface {
	GetTarget(url *url.URL) (Rule, error)
	Eval(io.Reader) error
}
