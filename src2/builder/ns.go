package builder

import (
	"errors"
	"fmt"
	"runtime"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/graph"
	"bldy.build/bldy/src/namespace"
)

var (
	ErrHostNotAvailable = errors.New("this compilation target is not compatible to run on this plan")
)

func nodeid(n *graph.Node) string {
	return fmt.Sprintf("%s-%s-bldy-%s-%x", n.ID, runtime.GOARCH, runtime.GOOS, n.Sum())
}

func (b *Builder) newnamespace(n *graph.Node, rt build.Runtime) (namespace.Namespace, error) {
	id := fmt.Sprintf("%x", n.Sum())
	return rt.NewNamespace(id)
}
