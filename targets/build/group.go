// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"path/filepath"

	"bldy.build/build"
)

type Group struct {
	Name         string   `group:"name"`
	Dependencies []string `group:"deps"`
	Exports      map[string]string
	Prefix       *string `group:"prefix" group:"prefix" build:"expand"`
}

func (g *Group) Hash() []byte {
	return []byte(g.Name)
}

func (g *Group) Build(c *build.Runner) error {
	return nil
}

func (g *Group) GetName() string {
	return g.Name
}

func (g *Group) GetDependencies() []string {
	return g.Dependencies
}
func (g *Group) Installs() map[string]string {
	if g.Prefix != nil {
		m := make(map[string]string)
		for dst, src := range g.Exports {
			m[filepath.Join(*g.Prefix, dst)] = src
		}
		return m
	}
	return g.Exports
}
