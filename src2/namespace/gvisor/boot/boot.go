package boot

import (
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"gvisor.dev/gvisor/pkg/control/server"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/p9"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/unet"
	"gvisor.dev/gvisor/runsc/fsgofer"
)

func Boot(args Args) (*kernel.Kernel, error) {
	l, _, err := createLoader()
	if err != nil {
		return nil, err
	}
	return l.k, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	if err := fsgofer.OpenProcSelfFD(); err != nil {
		panic(err)
	}
}

func testConfig() *Config {
	return &Config{
		RootDir:        "unused_root_dir",
		Network:        NetworkNone,
		DisableSeccomp: true,
		Platform:       "ptrace",
	}
}

// testSpec returns a simple spec that can be used in tests.
func testSpec() *specs.Spec {
	return &specs.Spec{
		// The host filesystem root is the sandbox root.
		Root: &specs.Root{
			Path:     "/",
			Readonly: true,
		},
		Process: &specs.Process{
			Args: []string{"/bin/true"},
		},
	}
}

// startGofer starts a new gofer routine serving 'root' path. It returns the
// sandbox side of the connection, and a function that when called will stop the
// gofer.
func startGofer(root string) (int, func(), error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return 0, nil, err
	}
	sandboxEnd, goferEnd := fds[0], fds[1]

	socket, err := unet.NewSocket(goferEnd)
	if err != nil {
		syscall.Close(sandboxEnd)
		syscall.Close(goferEnd)
		return 0, nil, fmt.Errorf("error creating server on FD %d: %v", goferEnd, err)
	}
	at, err := fsgofer.NewAttachPoint(root, fsgofer.Config{ROMount: true})
	if err != nil {
		return 0, nil, err
	}
	go func() {
		s := p9.NewServer(at)
		if err := s.Handle(socket); err != nil {
			log.Infof("Gofer is stopping. FD: %d, err: %v\n", goferEnd, err)
		}
	}()
	// Closing the gofer socket will stop the gofer and exit goroutine above.
	cleanup := func() {
		if err := socket.Close(); err != nil {
			log.Warningf("Error closing gofer socket: %v", err)
		}
	}
	return sandboxEnd, cleanup, nil
}

func createLoader() (*Loader, func(), error) {
	fd, err := server.CreateSocket(ControlSocketAddr(fmt.Sprintf("%010d", rand.Int())[:10]))
	if err != nil {
		return nil, nil, err
	}
	conf := testConfig()
	spec := testSpec()

	sandEnd, cleanup, err := startGofer(spec.Root.Path)
	if err != nil {
		return nil, nil, err
	}

	stdio := []int{int(os.Stdin.Fd()), int(os.Stdout.Fd()), int(os.Stderr.Fd())}
	args := Args{
		ID:           "foo",
		Spec:         spec,
		Conf:         conf,
		ControllerFD: fd,
		GoferFDs:     []int{sandEnd},
		StdioFDs:     stdio,
	}
	l, err := New(args)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return l, cleanup, nil
}
