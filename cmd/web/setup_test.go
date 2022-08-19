package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Do something first and run the test then exit.
	os.Exit(m.Run())
}

type MyHandler struct {
}

func (mh *MyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
