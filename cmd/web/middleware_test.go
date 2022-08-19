package main

import (
	"net/http"
	"testing"
)

func TestWriteToConsole(t *testing.T) {

	var myH MyHandler
	h := WriteToConsole(&myH)
	switch v := h.(type) {
	// do nothing when case is http.Handler
	case http.Handler:

	default:
		t.Errorf("Type is not http handler instead its type is %T", v)
	}

}

func TestNoSurf(t *testing.T) {

	var myH MyHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	// do nothing when case is http.Handler
	case http.Handler:

	default:
		t.Errorf("Type is not http handler instead its type is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {

	var myH MyHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	// do nothing when case is http.Handler
	case http.Handler:

	default:
		t.Errorf("Type is not http handler instead its type is %T", v)
	}
}
