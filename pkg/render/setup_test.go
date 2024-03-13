package render

import (
	"net/http"
	"os"
	"testing"

	"nicksrepo.com/nick/pkg/config"
)

var appTest config.App

type myWriter struct{}

func runALLtests(m *testing.M) {
	os.Exit(m.Run())
}

func (tw myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw myWriter) WriteHeader(i int) {}

func (tw myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
