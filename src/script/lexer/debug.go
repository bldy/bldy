package lexer

import (
	"runtime"
	"strings"
)

func (l *Lexer) Debug() {
	l.debug = !l.debug
}

func caller() (call string, file string, line int) {
	var caller uintptr
	caller, file, line, _ = runtime.Caller(2)
	name := strings.Split(runtime.FuncForPC(caller).Name(), ".")
	callName := name[len(name)-1]
	return callName, file, line
}
