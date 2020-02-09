package trace

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

var syzcalls = []uint64{
	unix.SYS_EXECVE,
	unix.SYS_CLONE,
	unix.SYS_FORK,
	unix.SYS_GETPID,
}

func allowed(kall uint64) bool {
	for _, k := range syzcalls {
		if kall == k {
			return true
		}
	}
	return false
}

type Tracer struct {
	pid         int
	children    []Tracer
	currentCall unix.PtraceRegsAmd64
}

func Trace(a string, args ...string) *Tracer {
	cmd := exec.Command(a, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned: %v\n", err)
	}

	p := &Tracer{pid: cmd.Process.Pid}
	p.trace()

	return p
}

func stringArgument(pid int, addr uintptr) string {
	var buffer [4096]byte
	n, err := syscall.PtracePeekData(pid, uintptr(addr), buffer[:])
	if err != nil {
		return ""
	}

	k := bytes.IndexByte(buffer[:n], 0)
	if k <= n {
		n = k
	}
	return string(buffer[:n])
}

type Execve struct {
	syscall unix.PtraceRegsAmd64
	pid     int
	file    string
	args    []string
}

func (e *Execve) String() string {
	return fmt.Sprintf(`execve	<%d>	%q`, e.pid, e.file)
}

type Clone struct {
	syscall unix.PtraceRegsAmd64
	pid     int
	newPid  int
}

func (c *Clone) String() string {
	return fmt.Sprintf(`clone	<%d>	newpid: %d`, c.pid, c.newPid)
}

type SysCall interface {
	String() string
}

func (p *Tracer) toTraceCall(r unix.PtraceRegsAmd64) SysCall {
	switch r.Orig_rax {
	case unix.SYS_EXECVE:
		return &Execve{
			pid:     p.pid,
			syscall: r,
			file:    stringArgument(p.pid, uintptr(r.Rdi)),
		}
	case unix.SYS_CLONE:
		return &Clone{
			pid:     p.pid,
			syscall: r,
			newPid:  int(r.Rax),
		}
	default:
		return nil
	}
}

func (p *Tracer) trace() {

	exit := true
	for {

		if err := unix.PtraceGetRegsAmd64(p.pid, &p.currentCall); err != nil {
			break
		}

		switch call := p.toTraceCall(p.currentCall).(type) {
		case *Execve:
			log.Println(call)
		case *Clone:
			if call.newPid > 0 {
				log.Println(call)
				_ = &Tracer{pid: call.newPid}
				//		log.Println(t.cmdline)
				if err := unix.PtraceAttach(call.newPid); err != nil {
					panic(err)
				}
			}
		}

		if err := syscall.PtraceSyscall(p.pid, 0); err != nil {
			panic(err)
		}
		if _, err := syscall.Wait4(p.pid, nil, 0, nil); err != nil {
			panic(err)
		}
		exit = !exit
	}
	return
}

func (t *Tracer) cmdline() string {
	statPath := fmt.Sprintf("/proc/%d/cmdline", t.pid)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return ""
	}

	// First, parse out the image name
	s := []string{""}
	for _, c := range dataBytes {
		if c != 0 {
			s[len(s)-1] += string(c)
		} else {
			s = append(s, "")
		}
	}
	if s[len(s)-1] == "" {
		s = s[:len(s)-1]
	}
	return fmt.Sprint(strings.Join(s, " "))
}
