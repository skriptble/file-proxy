package froxy

import (
	"errors"
	"io"
)

var NotFound = errors.New("File could not be found")

type Proxy interface {
	// RetrieveFile takes a name and the source and returns an io.ReadCloser
	// and an error. A nil ReadCloser or an error will result in a 404.
	RetrieveFile(name string, source string) (io.ReadCloser, error)
}

// ProxyBuilder is the interface used to construct an implementation of
// a Proxy.
type ProxyBuilder interface {
	Proxy
	// AddFileSource takes a FileSource and the source name and adds it to the proxy
	// returning an error if a problem occurs.
	AddFileSource(FileSource, string) error
}

// FileSource is the interface used by the Proxy to retrieve files.
type FileSource interface {
	// Open takes a file name and returns a ReadCloser and an error.
	Open(name string) (io.ReadCloser, error)
}
