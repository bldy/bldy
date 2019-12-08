package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/graph"
)

const (
	SCSSLOG = "success"
	FAILLOG = "fail"
)

func (b *Builder) buildpath(n *graph.Node) string {
	return filepath.Join(
		*b.config.Cache,
		nodeid(n),
	)
}

func (b *Builder) cached(n *graph.Node) bool {
	_, err := os.Lstat(b.buildpath(n))
	n.Cached = !os.IsNotExist(err)
	return n.Cached
}

func (b *Builder) builderror(n *graph.Node) error {
	nspath := b.buildpath(n)
	if file, err := os.Open(filepath.Join(nspath, FAILLOG)); err == nil {
		n.Status = build.Fail
		errString, _ := ioutil.ReadAll(file)
		return fmt.Errorf("%s", errString)
	} else if _, err := os.Lstat(filepath.Join(nspath, SCSSLOG)); err == nil {
		n.Status = build.Success
	}
	return nil
}

func (b *Builder) saveLog(n *graph.Node) {
	logName := "build.log"
	switch n.Status {
	case build.Success:
		logName = SCSSLOG
	case build.Fail:
		logName = FAILLOG
	}
	_ = logName
	/*	if logfile, err := os.Create(filepath.Join(b.buildpath(n), logName)); err != nil {
			l.Fatalf("error creating log for %s: %s", n.Task.Name(), err.Error())
		} else {
			if _, err := io.WriteString(logfile, n.Output); err != nil {
				l.Fatalf("error writing log for %s: %s", n.Task.Name(), err.Error())
			}
		}*/
}
