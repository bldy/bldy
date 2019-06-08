package gvisor

import (
	"context"
	"testing"
	"time"

	"gvisor.dev/gvisor/pkg/sentry/kernel"

	"bldy.build/bldy/src/namespace/gvisor/boot"
)

func TestBoot(t *testing.T) {
	k, err := boot.Boot(defaultArgs())
	k.Start()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	go func() {

		tg, id, err := k.CreateProcess(kernel.CreateProcessArgs{Argv: []string{"/bin/assdf"}})
		panic("SDF")
		if err != nil {
			t.Log("SDFS")
		}
		go k.StartProcess(tg)

		tg.WaitExited()
		cancel()
		_ = id
	}()
	_ = ctx
	select {
	case <-ctx.Done():
		panic(ctx.Err()) // prints "context deadline exceeded"
	}

}
