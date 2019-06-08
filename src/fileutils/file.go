package fileutils

import (
	"net/url"
	"os"
)

// ResolveFromWD will take string s and parse it as a URL
// will than resolve it's URL relative to the WorkDir
// if at anypoint it fails to
func ResolveFromWD(s string) (*url.URL, error) {
	// get the working directory
	wd, err := os.Getwd()
	base, err := url.Parse(wd + "/") // this is a directory so add a trailing stash at the end
	if err != nil {
		return nil, err
	}
	base.Scheme = "file"

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return base.ResolveReference(u), nil
}