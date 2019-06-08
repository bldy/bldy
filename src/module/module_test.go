package module

import (
	"fmt"
	"testing"

	"bldy.build/bldy/fileutils"
	_ "bldy.build/bldy/repository/file"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "hello",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := fileutils.ResolveFromWD(fmt.Sprintf(fmt.Sprintf("testdata/%s/test.bldy", test.name)))
			mod, err := New(u)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}
			t.Log(mod)
		})

	}
}