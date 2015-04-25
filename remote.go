package froxy

import (
	"io"
	"net/http"
	"net/url"
	"path"
)

// RemoteSource handles fetching files from a remote url destination.
type remote struct {
	href url.URL
}

// NewRemote creates a new FileSource that retrieves files from the given url.
func NewRemote(href url.URL) FileSource {
	r := remote{href: href}
	return r
}

// Open implements the Open method for the FileSource interface.
func (r remote) Open(name string) (io.ReadCloser, error) {
	href := r.href
	href.Path = path.Join(href.Path, name)
	// Make the request
	resp, err := http.Get(href.String())
	if err != nil {
		return nil, err
	}
	// Ensure we got a 2xx, if not return a 404
	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return nil, NotFound
	}
	// Get the body and just return that.
	return resp.Body, nil
}
