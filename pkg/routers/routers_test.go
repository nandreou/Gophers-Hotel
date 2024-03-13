package routers

import (
	"testing"

	"github.com/go-chi/chi"
	"nicksrepo.com/nick/pkg/middlewares"
)

func TestRoutes(t *testing.T) {
	middlewares.PassConfigToMidPkg(&app)
	mux := Routers()

	switch mux.(type) {
	case *chi.Mux:
	default:
		t.Error("Error")

	}

}
