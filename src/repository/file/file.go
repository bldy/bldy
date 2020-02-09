package file

import (
	"io"
	"net/url"
	"os"

	"bldy.build/bldy/repository"
)

func init() {
	repository.RegisterScheme("file", Dial)
}

func Dial(u *url.URL) (repository.Repository, error) {
	return &fileRepo{base: &(*u)}, nil
}

type fileRepo struct{ base *url.URL }

func (r *fileRepo) Open(name string) (io.Reader, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}
	n := r.base.ResolveReference(u)
	return os.Open(n.Path)
}
