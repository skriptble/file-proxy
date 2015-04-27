package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skriptble/froxy"
)

func TestProxy(t *testing.T) {
	wrapper := proxyWrapper{}
	// Should respond with error 405 to non-GET request
	req, err := http.NewRequest("POST", "http://example.com/", nil)
	if err != nil {
		t.Errorf("An unexpected error occured %v", err)
	}
	w := httptest.NewRecorder()
	wrapper.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Error("Should respond with error 405 to non-GET request")
		t.Errorf("Wanted %v, got %v", http.StatusMethodNotAllowed, w.Code)
	}
	// Should respond with error when path is not of the correct length
	req, err = http.NewRequest("GET", "http://example.com/local", nil)
	if err != nil {
		t.Errorf("An unexpected error occured %v", err)
	}
	w = httptest.NewRecorder()
	wrapper.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Should respond with error when path is not of the correct length")
		t.Errorf("Wanted %v, got %v", http.StatusNotFound, w.Code)
	}
	want := froxy.NotFound.Error()
	got := w.Body.String()
	if got != want {
		t.Error("Should respond with error when path is not of the correct length")
		t.Errorf("Wanted %v, got %v", want, got)
	}
	// Should handle error from a proxy
	proxy := froxy.NewProxy()
	proxy.AddFileSource(dummyFileSource{}, "dummy")
	wrapper = proxyWrapper{p: proxy}
	req, err = http.NewRequest("GET", "http://example.com/dummy/foo-bar", nil)
	if err != nil {
		t.Errorf("An unexpected error occured %v", err)
	}
	w = httptest.NewRecorder()
	wrapper.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Should handle error from a proxy")
		t.Errorf("Wanted %v, got %v", http.StatusNotFound, w.Code)
	}
	want = ErrDummy.Error()
	got = w.Body.String()
	if got != want {
		t.Error("Should handle error from a proxy")
		t.Errorf("Wanted %v, got %v", want, got)
	}

	// Should return a file retrieved from a proxy
	req, err = http.NewRequest("GET", "http://example.com/dummy/dummy", nil)
	if err != nil {
		t.Errorf("An unexpected error occured %v", err)
	}
	w = httptest.NewRecorder()
	wrapper.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Should return a file retrieved from a proxy")
		t.Errorf("Wanted %v, got %v", http.StatusOK, w.Code)
	}
	want = "dummy file data foo bar baz qux quux"
	got = w.Body.String()
	if got != want {
		t.Error("Should return a file retrieved from a proxy")
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

var ErrDummy = errors.New("dummy error occured")

type dummyFileSource struct {
}

func (dfs dummyFileSource) Open(name string) (io.ReadCloser, error) {
	if name == "dummy" {
		reader := bytes.NewReader([]byte("dummy file data foo bar baz qux quux"))
		readcloser := ioutil.NopCloser(reader)
		return readcloser, nil
	}
	return nil, ErrDummy
}
