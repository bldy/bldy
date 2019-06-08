// Copyright 2015-2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graph parses and generates build graphs
package graph

import (
	"fmt"
	"log"
	"os"

	"bldy.build/bldy/src/url"

	"bldy.build/bldy/src/build"
)

var (
	l = log.New(os.Stdout, "graph: ", 0)
)

type ResolverFunc func(u *url.URL) (build.Task, error)

// Graph represents a build graph
type Graph struct {
	Root  *Node
	Nodes map[string]*Node

	resolve ResolverFunc
}

// New returns a new build graph relatvie to the working directory
func New(ctx build.Context, u *url.URL, resolver ResolverFunc) (*Graph, error) {
	g := Graph{
		Nodes:   make(map[string]*Node),
		resolve: resolver,
	}
	var err error
	if g.Root, err = g.getTask(u); err != nil {
		return nil, fmt.Errorf("graph new: %v", err)
	}
	g.Root.IsRoot = true
	return &g, nil
}

func (g *Graph) getTask(u *url.URL) (*Node, error) {
	if gnode, ok := g.Nodes[u.String()]; ok {
		return gnode, nil
	}

	t, err := g.resolve(u)
	if err != nil {
		return nil, fmt.Errorf("get task: %v", err)
	}

	node := NewNode(u, t)

	var deps []build.Task

	for _, d := range node.Task.Dependencies() {
		c, err := g.getTask(d)
		if err != nil {
			return nil, fmt.Errorf("get dependencies: %v", err)
		}

		node.WG.Add(1)

		deps = append(deps, c.Task)

		node.Children[d.String()] = c
		c.Parents[u.String()] = &node
	}

	g.Nodes[u.String()] = &node
	if t.Name() == u.String() {
		return &node, nil
	}
	return nil, fmt.Errorf("graph: target name %q and url target %q don't match", t.Name(), u.String())

}
