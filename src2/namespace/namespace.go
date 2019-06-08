package namespace

import (
	"os"

	_ "gvisor.dev/gvisor/pkg/sentry/platform/ptrace"
)

type MountFlag int

const (
	Replace MountFlag = iota
	Before
	After
)

func New() *Namespace {
	return &Namespace{}
}

type bind struct {
	new, old string
	flags    MountFlag
}
type mount struct {
	new, old string
	flags    MountFlag
}

type Namespace struct {
	binds  []bind
	mounts []mount
}

// Bind takes the portion of the existing name space visible at new,
// either a file or a directory, and makes it also visible at old.
// For example,
//
// 	bind("1995/0301/sys/include", "/sys/include", REPLACE)
//
func (ns *Namespace) Bind(new, old string, flags MountFlag) {
	ns.binds = append(ns.binds, bind{new, old, flags})
}

func (ns *Namespace) Mount(new, old string, flags MountFlag) {
	ns.mounts = append(ns.mounts, mount{new, old, flags})
}

func (ns *Namespace) Open(name string) (*os.File, error) {
	return nil, nil
}
func (ns *Namespace) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}
func (ns *Namespace) Create(name string) (*os.File, error) {
	return nil, nil
}

func Run() {

}