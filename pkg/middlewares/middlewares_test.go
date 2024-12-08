package middlewares

import (
	"net/http"
	"testing"

	"github.com/justinas/nosurf"
)

func TestCsfrf(t *testing.T) {
	PassConfigToMidPkg(&app2)
	csrf := CsrfMiddleWare(testHandler)

	switch csrf.(type) {
	case *nosurf.CSRFHandler:
	default:
		t.Error("Test Failed Function did not return a csrfToken")
	}
}

func TestSessionLoader(t *testing.T) {
	session := SessionLoader(testHandler)

	switch session.(type) {
	case http.HandlerFunc:
	default:
		t.Error("Test Failed Function did not return a csrfToken")
	}
}

func TestSetCookies(t *testing.T) {
	cookies := SetCookies(testHandler)

	switch cookies.(type) {
	case http.HandlerFunc:
	default:
		t.Error("Test Failed Function did not return a csrfToken")
	}
}
