package ziggy

import (
	"context"
	"os"
	"path"
	"testing"

	"bldy.build/bldy/src/build"
	"bldy.build/bldy/src/namespace/gvisor"
	"bldy.build/bldy/src/url"
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

func TestDoesntExist(t *testing.T) {
	u, err := url.Parse("empty#empty")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	z := New(c, &gvisor.Runtime{})
	r, err := z.GetTask(u)
	if err == nil || r != nil {
		t.Log("did not expwct a target")
		t.Fail()
	}
}

/*
func TestEval(t *testing.T) {
	wd, _ := os.Getwd()
	tests := []struct {
		name string
		url  string
		wd   string
		err  error
	}{
		{
			name: "empty",
			url:  "empty",
			wd:   path.Join(wd, "testdata"),
			err:  nil,
		},
		{
			name: "context_tester",
			url:  "ctx#context_tester",
			wd:   path.Join(wd, "testdata"),
			err:  nil,
		},
		{
			name: "run",
			url:  "run#run_test",
			wd:   path.Join(wd, "testdata"),
			err:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vm := New(c, &gvisor.Runtime{})
			u, _ := url.Parse(test.url)
			target, err := vm.GetTask(u)
			if err != test.err {
				t.Log(err)
				t.Fail()
				return
			}
			if target == nil {
				t.Fail()
				return
			}
			if target == nil {
				t.Fail()
				return
			}
		})
	}
}
*/