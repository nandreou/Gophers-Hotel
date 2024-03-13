package handlerfuntions

import (
	"encoding/gob"
	"net/http"
	"net/url"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"nicksrepo.com/nick/pkg/config"
	"nicksrepo.com/nick/pkg/database"
	"nicksrepo.com/nick/pkg/middlewares"
	"nicksrepo.com/nick/pkg/render"
)

var app config.App
var db database.DB

var keyValuePairBookNow = [][]struct {
	key   string
	value string
}{
	{
		{"first_name", "nick"},
		{"last_name", "nick"},
		{"email", "nick@gmail.com"},
		{"phone", "245845687"},
	},
	{
		{"first_name", "nick"},
		{"last_name", "nick"},
		{"email", "nick@gmail.com"},
		{"phone", "24354"},
	},
}

var reqBodyPost = []struct {
	path     string
	keyValue [][]struct {
		key   string
		value string
	}
	expected []int
}{
	{"/book-now", keyValuePairBookNow, []int{http.StatusSeeOther, http.StatusOK}},
}

var reqBodyGet = []struct {
	path     string
	expected int
}{
	{"/", http.StatusOK},
	{"/about", http.StatusOK},
	{"/generals-quarters", http.StatusOK},
	{"/majors-suite", http.StatusOK},
	{"/book-now", http.StatusOK},
	{"/reservation-summary", http.StatusOK},
	{"/contact", http.StatusOK},
}

func setUpRoutes() http.Handler {
	var cacheError error

	gob.Register(url.Values{})

	//Cache the Templates in Servers Memory
	app.Template, cacheError = render.CreateTmplCache()

	if cacheError != nil {
		return nil
	}

	app.Prd = false

	//Create Session Manager and session configs
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Prd

	app.Session = session

	//Pass App config In Pakages
	render.NewCache(&app)
	middlewares.PassConfigToMidPkg(&app)

	//Create Repository
	NewRepository(&app, &db)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(middlewares.CsrfMiddleWare)
	mux.Use(middlewares.SetCookies)
	mux.Use(middlewares.SessionLoader)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.GeneralQuarters)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/book-now", Repo.BookNow)
	mux.Post("/book-now", Repo.BookNow)
	mux.Get("/reservation-summary", Repo.ReservationSumary)

	mux.Get("/contact", Repo.Contact)

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("../static/"))))

	return mux
}
