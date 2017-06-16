package cc

import (
	"crypto/sha1"
	"fmt"

	"io"
	"strings"

	"path/filepath"

	"bldy.build/build"
	"bldy.build/build/racy"
)

type CBin struct {
	Name            string        `cxx_binary:"name" cc_binary:"name"`
	Sources         []string      `cxx_binary:"srcs" cc_binary:"srcs" build:"path" ext:".c,.S,.cpp,.cc"`
	Dependencies    []string      `cxx_binary:"deps" cc_binary:"deps"`
	Includes        Includes      `cxx_binary:"headers" cc_binary:"includes" build:"path" ext:".h,.c,.S"`
	Headers         []string      `cxx_binary:"exported_headers" cc_binary:"hdrs" build:"path" ext:".h,.hh,.hpp"`
	CompilerOptions CompilerFlags `cxx_binary:"compiler_flags" cc_binary:"copts"`
	LinkerOptions   []string      `cxx_binary:"linker_flags" cc_binary:"linkopts"`
	LinkerFile      string        `cxx_binary:"ld" cc_binary:"ld" build:"path"`
	Static          bool          `cxx_binary:"linkstatic" cc_binary:"linkstatic"`
	Strip           bool          `cxx_binary:"strip" cc_binary:"strip"`
	AlwaysLink      bool          `cxx_binary:"alwayslink" cc_binary:"alwayslink"`
	Install         *string       `cxx_binary:"install" cc_binary:"install"`
}

func split(s string, c string) string {
	i := strings.Index(s, c)
	if i < 0 {
		return s
	}

	return s[i+1:]
}
func (cb *CBin) Hash() []byte {
	h := sha1.New()
	io.WriteString(h, CCVersion)
	io.WriteString(h, cb.Name)
	racy.HashStrings(h, cb.CompilerOptions)
	racy.HashStrings(h, cb.LinkerOptions)
	return racy.XOR(
		h.Sum(nil),
		racy.HashFilesForExt([]string(cb.Includes), ".h"),
		racy.HashFilesForExt(cb.Sources, ".c"),
		racy.HashFilesForExt(cb.Sources, ".S"),
	)
}

func (cb *CBin) Build(e *build.Executor) error {
	params := []string{"-c"}
	params = append(params, cb.CompilerOptions...)
	params = append(params, cb.Sources...)

	params = append(params, cb.Includes.Includes()...)

	if err := e.Exec(Compiler(), CCENV, params); err != nil {
		return err
	}

	ldparams := []string{"-o", cb.Name}

	// This is done under the assumption that each src file put in this thing
	// here will comeout as a .o file
	for _, f := range cb.Sources {
		_, fname := filepath.Split(f)
		ldparams = append(ldparams, fmt.Sprintf("%s.o", fname[:strings.LastIndex(fname, ".")]))
	}

	ldparams = append(ldparams, cb.LinkerOptions...)
	if cb.LinkerFile != "" {
		ldparams = append(ldparams, cb.LinkerFile)
	}
	haslib := false
	for _, dep := range cb.Dependencies {
		d := split(dep, ":")
		if len(d) < 3 {
			continue
		}
		if d[:3] == "lib" {
			if cb.AlwaysLink {
				ldparams = append(ldparams, fmt.Sprintf("%s.a", d))
			} else {
				if !haslib {
					ldparams = append(ldparams, "-L", "lib")
					haslib = true
				}
				ldparams = append(ldparams, fmt.Sprintf("-l%s", d[3:]))
			}
		}

		// kernel specific
		if len(d) < 4 {
			continue
		}
		if d[:4] == "klib" {
			ldparams = append(ldparams, fmt.Sprintf("%s.a", d))
		}
	}

	if err := e.Exec(Linker(), CCENV, ldparams); err != nil {
		return err
	}
	if cb.Strip {
		sparams := []string{"-o", cb.Name, cb.Name}
		if err := e.Exec(Stripper(), nil, sparams); err != nil {
			return err
		}
	}
	return nil
}

func (cb *CBin) Installs() map[string]string {
	exports := make(map[string]string)
	if cb.Install != nil {
		exports[*cb.Install] = cb.Name
	} else {
		exports[filepath.Join("bin", cb.Name)] = cb.Name
	}
	return exports
}

func (cb *CBin) GetName() string {
	return cb.Name
}
func (cb *CBin) GetDependencies() []string {
	return cb.Dependencies
}
