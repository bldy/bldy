package integration

import (
	"context"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/graph"
	"bldy.build/bldy/src/url"

	_ "bldy.build/bldy/src/ziggy/local"
)

type testContext struct {
	context.Context
	base url.URL
}

func (c *testContext) Getbase() (base url.URL) { return c.base }
func (c *testContext) WithBase(base *url.URL) build.Context {
	ctx := *c
	ctx.base = *base
	return &ctx
}

var c = func() build.Context {
	wd, _ := os.Getwd()
	u, _ := url.Parse("file://" + path.Join(wd, "testdata"))
	return &testContext{context.Background(), *u}
}()

var tests = []struct {
	name string
	url  string
	wd   string
	err  error
}{

	{
		name: "run",
		url:  "rust/rust.bldy#rustbin",
		err:  nil,
	},
}

type testNotifier struct {
	t *testing.T
}

func (t *testNotifier) Update(n *graph.Node) {
	switch n.Status {
	default:
		t.t.Logf("%s %s ", n.Status, n.ID)
	}

}

func (t *testNotifier) Error(err error) {
	t.t.Fail()
	t.t.Logf("error: %+v\n", err)
}

func (t *testNotifier) Done(d time.Duration) {
	t.t.Logf("Finished building in %s\n", d)

}

func TestBuild(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log("bldy", "build", test.url)
			wd, _ := os.Getwd()
			cmd := exec.Command("bldy", "build", test.url)
			cmd.Dir = path.Join(wd, "testdata")
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				os.Stdout.Write(stdoutStderr)
				t.Log(err)
				t.Fail()
			}
		})
	}
}