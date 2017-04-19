package harvey

import (
	"crypto/sha1"

	"io"

	"bldy.build/build"
)

type Move struct {
	Name         string            `move:"name"`
	Dependencies []string          `move:"deps"`
	Exports      map[string]string `move:"installs" build:"expand"`
}

func (m *Move) Hash() []byte {
	h := sha1.New()
	io.WriteString(h, m.Name)
	return h.Sum(nil)
}

func (m *Move) Build(c *build.Runner) error {
	return nil
}

func (m *Move) Installs() map[string]string {
	return m.Exports
}

func (m *Move) GetName() string {
	return m.Name
}
func (m *Move) GetDependencies() []string {
	return m.Dependencies
}
