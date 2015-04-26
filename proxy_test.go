package froxy

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestProxy(t *testing.T) {
	// Should be able to add a FileSource to the proxy
	p := NewProxy()
	dummy := dummyFileSource{}
	err := p.AddFileSource(dummy, "dummy")
	if err != nil {
		t.Error("Should be able to add a FileSource to the proxy")
		t.Errorf("Wanted nil, got %v", err)
	}
	// Should be able to open a file from the FileSource
	resp, err := p.RetrieveFile("dummy", "dummy")
	if err != nil {
		t.Error("Should be able to open a file from the FileSource")
		t.Errorf("Wanted nil, got %v", err)
	}
	want := []byte("dummy file data foo bar baz qux quux")

	got, err := ioutil.ReadAll(resp)
	if err != nil {
		t.Errorf("An unexpected error occured %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Error("Should be able to open a file from the FileSource")
		t.Errorf("Wanted %v, got %v", want, got)
	}
	// Should properly handle errors from FileSource
	_, err = p.RetrieveFile("", "dummy")
	if !reflect.DeepEqual(ErrDummy, err) {
		t.Error("Should properly handle errors from FileSource")
		t.Errorf("Wanted %v, got %v", ErrDummy, err)
	}
	// Should return NotFound error when the FileSource does not exist for the
	// specified source.
	_, err = p.RetrieveFile("", "doesnt-exist")
	if !reflect.DeepEqual(NotFound, err) {
		t.Error("Should return NotFound error when the FileSource does not exist for the specified source.")
		t.Errorf("Wanted %v, got %v", NotFound, err)
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
