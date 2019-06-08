package url

import (
	"fmt"
	"net/url"
	"path"
)

type URL struct {
	url.URL
}

func Parse(s string) (*URL, error) {
	u, err := url.Parse(s)
	return &URL{*u}, err
}

func (u *URL) IsAbs() bool { return path.IsAbs(u.String()) }

func (u *URL) Append(b *URL) (*URL, error) {
	a := u
	if b.IsAbs() && a.Scheme != b.Scheme {
		return nil, fmt.Errorf("url: appended strings must be of the same scheme")
	}
	aa, bb := *a, *b
	aa.Scheme = ""
	bb.Scheme = ""
	if x, err := Parse(path.Join(aa.String(), bb.String())); err == nil {
		x.Scheme = u.Scheme
		return x, nil
	} else {
		return nil, err
	}
}