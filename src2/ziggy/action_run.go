package ziggy

import (
	"fmt"
	"os"

	"bldy.build/bldy/src/executor"
	"bldy.build/bldy/src/racy"
)

type run struct {
	Outputs               []string          `ziggy:"outputs"`                // List of the output files of the action.
	Files                 []string          `ziggy:"files"`                  // List of the input files of the action.
	Executable            string            `ziggy:"executable"`             // The executable file to be called by the action.
	Arguments             []string          `ziggy:"arguments"`              // Command line arguments of the action. Must be a list of strings or actions.args() objects.
	Mnemonic              string            `ziggy:"mnemonic"`               // A one-word description of the action, for example, CppCompile or GoLink.
	ProgressMessage       string            `ziggy:"progress_message"`       // Progress message to show to the user during the build, for example, "Compiling foo.cc to create foo.o".
	UseDefaultShellEnv    bool              `ziggy:"use_default_shell_env"`  // Whether the action should use the built in shell environment or not.
	Env                   map[string]string `ziggy:"env"`                    // Sets the dictionary of environment variables.
	ExecutionRequirements map[string]string `ziggy:"execution_requirements"` // Information for scheduling the action. See tags for useful keys.
}

func (r *run) Sum() []byte {
	h := racy.New()
	h.HashStrings(r.Executable)
	h.HashStrings(r.Arguments...)
	return h.Sum(nil)
}

func (r *run) Do(e *executor.Executor) error {
	env := os.Environ()
	if !r.UseDefaultShellEnv {
		env = []string{}
	}
	for k, v := range r.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	e.Printf(r.ProgressMessage)
	return e.Exec(r.Executable, env, r.Arguments)
}
