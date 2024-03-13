package middlewares

import (
	"net/http"
	"os"
	"testing"

	"nicksrepo.com/nick/pkg/config"
)

var testHandler http.Handler
var app2 config.App

func MiddleTests(m testing.M) {
	os.Exit(m.Run())
}
