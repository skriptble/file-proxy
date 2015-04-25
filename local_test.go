package froxy

import (
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestLocal(t *testing.T) {
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
	local := Dir(".")

	// Should be able to open a file that exists
	rc, err := local.Open(fileName)
	if err != nil {
		t.Error("Should be able to open a file that eixsts")
		t.Errorf("Wanted nil, got %v", err)
	}
	got, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Error("Should be able to open a file that eixsts")
		t.Errorf("Wanted nil, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("Should be able to open a file that eixsts")
		t.Errorf("Wanted %v, got %v", want, got)
	}

	// Should not be able to open a file that doesn't exist
	_, err = local.Open("this-should-fail-hard")
	if err == nil {
		t.Error("Should not be able to open a file that doesn't exist")
	}
}
