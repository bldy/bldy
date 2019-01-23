package builder

import (
	"errors"
	"fmt"
	"runtime"

	"bldy.build/build/graph"
	"bldy.build/build/namespace"
)

var (
	ErrHostNotAvailable = errors.New("this compilation target is not compatible to run on this plan")
)

func nodeid(n *graph.Node) string {
	return fmt.Sprintf("%s-%s-bldy-%s-%x", n.Target.Name(), runtime.GOARCH, runtime.GOOS, n.HashNode())
}
func (b *Builder) newnamespace(n *graph.Node) (namespace.Namespace, error) {

	return nil, ErrHostNotAvailable
}
