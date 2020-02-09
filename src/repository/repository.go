package repository

import (
	"fmt"
	"io"
	"net/url"
	"sync"
)

type Repository interface {
	Open(name string) (io.Reader, error)
}

var (
	mu       sync.Mutex
	diallers = map[string]func(*url.URL) (Repository, error){}
)

// RegisterScheme associates a dialler with a URL scheme.
func RegisterScheme(scheme string, dial func(*url.URL) (Repository, error)) {
	mu.Lock()
	diallers[scheme] = dial
	mu.Unlock()
}

// Dial attempts to connect to the repository named by the given URL.
// The URL's scheme must be registered with RegisterScheme.
func Dial(u *url.URL) (Repository, error) {
	mu.Lock()
	dial := diallers[u.Scheme]
	mu.Unlock()
	if dial == nil {
		return nil, fmt.Errorf("dial: unknown scheme %q", u.Scheme)
	}
	return dial(u)
}