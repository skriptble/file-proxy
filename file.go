package froxy

import (
	"io"
	"net/http"
)

// Dir is just a wrapper around http.Dir.
type Dir string

// Open implements the Open method on the FileSource interface. This just wraps
// http.Dir's open method.
func (d Dir) Open(name string) (io.ReadCloser, error) {
	dir := http.Dir(d)
	return dir.Open(name)
}
