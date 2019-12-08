// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package build defines build target and build context structures
package build

import (
	"context"

	"bldy.build/bldy/src/executor"
	"bldy.build/bldy/src/namespace"
	"bldy.build/bldy/src/url"
)

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

// VM seperate the parsing and evauluating targets logic from rest of bldy
// so we can implement and use new grammars like jsonnet or go it self.
type Store interface {
	GetTask(url *url.URL) (Task, error)
}

// Runtime defines a runtime for the builds.
// This context is scoped to the specifics of the execution environment.
type Runtime interface {
	OS() string
	Arch() string

	Printf(formmat string, v ...interface{})
	NewNamespace(id string) (namespace.Namespace, error)
}

// Task defines is the execution specific Rule.
type Task interface {
	Name() string
	Sum() []byte
	Run(e *executor.Executor) error
	Dependencies() []*url.URL
	Outputs() []string
}

type Context interface {
	context.Context

	Getbase() (base url.URL)
	WithBase(u *url.URL) Context
}