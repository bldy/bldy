// Copyright 2017 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package project // import "bldy.build/build/project"

import (
	"log"
	"os"
)

var (
	l = log.New(os.Stdout, "project", 0)
)

// Getenv returns the envinroment variable. It looks for the envinroment
// variable in the following order. It checks if the current shell session has
// an envinroment variable, checks if it's set in the OS specific section in
// the .build file, and checks it for common in the .build config file.
func Getenv(s string) string {
	return os.Getenv(s)
}
