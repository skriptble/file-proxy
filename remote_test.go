package froxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestRemote(t *testing.T) {
	fileName := strconv.Itoa(int(time.Now().UnixNano()))
	file, err := os.Create(fileName)
	if err != nil {
		t.Errorf("An unexpected error occured: %v", err)
	}
	defer file.Close()
	defer os.Remove(fileName)
	want := []byte("Foo Bar Baz Qux Quux")
	_, err = file.Write(want)
	if err != nil {
		t.Errorf("An unexpected error occured: %v", err)
	}

	// Should be able to retrieve remote file
	ts := httptest.NewServer(http.FileServer(http.Dir(".")))
	defer ts.Close()

	href, err := url.Parse(ts.URL)
	if err != nil {
		t.Errorf("An unexpected error occured: %v", err)
	}
	rmt := NewRemote(*href)
	body, err := rmt.Open(fileName)
	if err != nil {
		t.Errorf("An unexpected error occured: %v", err)
	}
	got, err := ioutil.ReadAll(body)
	if err != nil {
		t.Errorf("An unexpected error occured: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("Should be able to retrieve remote file")
		t.Errorf("Wanted %v, got %v", want, got)
	}
	// Should return error when status code is 404
	_, err = rmt.Open("this-doesnt-exist-nor-should-it")
	if err == nil {
		t.Error("Should return error when status code is 404")
		t.Errorf("Wanted %v, got nil", NotFound)
	}
	// Should return an error when http.Get encounters an error
	rmt = NewRemote(url.URL{})
	_, err = rmt.Open("failure-of-epic-proportion")
	if err == nil {
		t.Error("Should return an error when http.Get encounters an error")
		t.Errorf("Wanted an error, got nil")
	}
}
