package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	handlerfuntions "nicksrepo.com/nick/pkg/handlerFuntions"
	"nicksrepo.com/nick/pkg/middlewares"
)

func Routers() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middlewares.CsrfMiddleWare)
	mux.Use(middlewares.SetCookies)
	mux.Use(middlewares.SessionLoader)

	mux.Get("/", handlerfuntions.Repo.Home)
	mux.Get("/about", handlerfuntions.Repo.About)

	mux.Get("/generals-quarters", handlerfuntions.Repo.GeneralQuarters)
	mux.Get("/majors-suite", handlerfuntions.Repo.Majors)

	mux.Get("/search-availability", handlerfuntions.Repo.SearchAvailability)
	mux.Post("/search-availability", handlerfuntions.Repo.SearchAvailability)

	mux.Post("/search-availability-by-room", handlerfuntions.Repo.SearchAvailabilityByRoom)

	mux.Get("/book-now/{id}", handlerfuntions.Repo.BookNow)
	mux.Post("/book-now/{id}", handlerfuntions.Repo.BookNow)

	mux.Get("/reservation-summary", handlerfuntions.Repo.ReservationSumary)

	mux.Get("/contact", handlerfuntions.Repo.Contact)

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("../static/"))))

	return mux
}
