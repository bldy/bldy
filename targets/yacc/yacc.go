// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"bytes"
	"crypto/sha1"

	"io"
	"os/exec"
	"path/filepath"

	"log"

	"strings"

	"fmt"

	"bldy.build/build"
	"bldy.build/build/internal"
	"bldy.build/build/racy"
)

var YaccVersion = ""

type Yacc struct {
	Name           string   `yacc:"name"`
	Sources        []string `yacc:"srcs" build:"path"`
	Exports        []string `yacc:"exports"`
	ExporedHeaders []string `yacc:"hdrs"`
	Dependencies   []string `yacc:"deps"`
	YaccOptions    []string `yacc:"yaccopts"`
	Source         string
	buf            bytes.Buffer
}

func init() {
	if out, err := exec.Command("yacc", "--version").Output(); err != nil {
		YaccVersion = "deadbeef"
	} else {
		YaccVersion = strings.TrimSpace(string(out))
	}

	if err := internal.Register("yacc", Yacc{}); err != nil {
		log.Fatal(err)
	}
}
func (y *Yacc) Hash() []byte {
	h := sha1.New()
	io.WriteString(h, YaccVersion)
	io.WriteString(h, y.Name)
	racy.HashFiles(h, y.Sources)
	racy.HashStrings(h, y.YaccOptions)
	return h.Sum(nil)
}

func (y *Yacc) Build(c *build.Runner) error {

	params := []string{}
	params = append(params, y.YaccOptions...)
	params = append(params, y.Sources...)

	c.Println(strings.Join(append([]string{"yacc"}, params...), " "))

	if err := c.Exec("yacc", nil, params); err != nil {
		c.Println(err.Error())
		return fmt.Errorf(y.buf.String())
	}

	return nil
}
func (y *Yacc) Installs() map[string]string {
	installs := make(map[string]string)
	for _, e := range y.Exports {
		installs[e] = e
	}
	for _, e := range y.ExporedHeaders {
		installs[filepath.Join("include", e)] = e
	}
	return installs
}
func (y *Yacc) GetName() string {
	return y.Name
}

func (y *Yacc) GetDependencies() []string {
	return y.Dependencies
}
func (y *Yacc) GetSource() string {
	return y.Source
}

func (y *Yacc) Reader() io.Reader {
	return &y.buf
}
