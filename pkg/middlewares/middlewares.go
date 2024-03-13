package middlewares

import (
	"net/http"

	"github.com/justinas/nosurf"
	"nicksrepo.com/nick/pkg/config"
)

var app *config.App

func PassConfigToMidPkg(a *config.App) {
	app = a
}

func CsrfMiddleWare(next http.Handler) http.Handler {
	csrf_token := nosurf.New(next)

	csrf_token.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.Prd,
		SameSite: http.SameSiteLaxMode,
	})

	return csrf_token
}

func SessionLoader(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

// Custom Middleware for coockie set
func SetCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "myown",
			Value: "gsadrgserg54",
		}
		http.SetCookie(w, cookie)

		cookie2 := &http.Cookie{
			Name:  "myowncoockie",
			Value: "asdfa4fafa4tge5g",
		}

		w.Header().Add("Cook", cookie2.String())
		next.ServeHTTP(w, r)
	})

}
