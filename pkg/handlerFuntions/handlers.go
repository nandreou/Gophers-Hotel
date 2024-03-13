package handlerfuntions

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"nicksrepo.com/nick/pkg/config"
	"nicksrepo.com/nick/pkg/database"
	"nicksrepo.com/nick/pkg/forms"
	"nicksrepo.com/nick/pkg/models"
	"nicksrepo.com/nick/pkg/render"
)

type Repository struct {
	App *config.App
	DB  database.DbHandler
}

var Repo Repository

func NewRepository(a *config.App, db *database.DB) {
	Repo.App = a
	Repo.DB = db
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	repo.App.Session.Put(r.Context(), "rem_ip", r.RemoteAddr)

	//Django Style
	data := &struct {
		Name    string
		SurName string
		Age     int
	}{
		"Nick",
		"Andreou",
		29,
	}

	render.Render(w, "home.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/ HTTP: 200")
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		//Django Style
		data := &struct {
			Name      string
			SurName   string
			Age       int
			CSRFtoken string
		}{
			"Nick",
			"Andreou",
			29,
			nosurf.Token(r),
		}

		render.Render(w, "about.page.tmpl", data)
		log.Println("GET:", r.RemoteAddr, "/about HTTP: 200")
	}
}

func (repo *Repository) GeneralQuarters(w http.ResponseWriter, r *http.Request) {

	//Django Style
	data := &struct {
		Name      string
		SurName   string
		Age       int
		CSRFtoken string
	}{
		"Nick",
		"Andreou",
		29,
		nosurf.Token(r),
	}

	render.Render(w, "generals.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/generals-quarters HTTP: 200")
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	//Django Style
	data := &struct {
		Name      string
		SurName   string
		Age       int
		CSRFtoken string
	}{
		"Nick",
		"Andreou",
		29,
		nosurf.Token(r),
	}

	render.Render(w, "majors.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/majors-suite HTTP: 200")
}

func (repo *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	if r.Method == "POST" {
		rooms, err := repo.DB.AvailabilitySearch(start, end)
		if err != nil {
			log.Println(err)
		}

		data := &struct {
			Rooms     []*models.SearchAvailabilityModel
			StartDate string
			EndDate   string
			CSRFtoken string
		}{
			rooms,
			start,
			end,
			nosurf.Token(r),
		}

		render.Render(w, "availability.page.tmpl", data)
		return
	}
	//Django Style
	data := &struct {
		CSRFtoken string
	}{
		nosurf.Token(r),
	}

	render.Render(w, "search-availability.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/search-availability HTTP: 200")
}

func (repo *Repository) BookNow(w http.ResponseWriter, r *http.Request) {

	/*IN THIS FUNCTION URL QUERIES ARE USED FOR START AND END DATES
	THE BEST PRACTICE HERE IS TO USE THE SESSION MECHANISM WITH A UNIQUE REQUEST ID FOR EVERY SO THE START AND END DATE ARE STORED.*/

	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error on writing data Try Again Later")
			return
		}

		form := forms.NewForm(r.Form)

		//SimpleValidateTest Form Validation
		for field := range r.Form {
			form.Has(field)
		}

		dateFormat := "2006-01-02"

		startDate, err := time.Parse(dateFormat, r.URL.Query().Get("startdate"))

		if err != nil {
			http.Error(w, "Failed to Make Reservation Please use 2024-01-02 StartDate Format", http.StatusBadRequest)
			log.Println("Failed to Write in StartDate Database", startDate, "\n", err)
			return
		}

		endDate, err := time.Parse(dateFormat, r.URL.Query().Get("enddate"))

		if err != nil {
			http.Error(w, "Failed to Make Reservation Please use 2024-01-02 EndDate Format", http.StatusBadRequest)
			log.Println("Failed to Write in EndDate Database", endDate, "\n", err)
			return
		}

		//Last check of availability
		id := chi.URLParam(r, "id")

		form.ValidForm, _ = repo.DB.LastAvailabilitySearch(id, startDate, endDate)

		//Check Form Validation
		if !form.IsValid() {
			log.Println("POST:", r.RemoteAddr, "/book-now HTTP: 400")
			http.Redirect(w, r, "/book-now/"+id+"?startdate="+r.URL.Query().Get("startdate")+"&&enddate="+r.URL.Query().Get("enddate"), http.StatusSeeOther)
			return
		}

		reservation := &models.Reservation{
			RoomId:    1,
			FirstName: r.Form.Get("first_name"),
			LastName:  r.Form.Get("last_name"),
			Email:     r.Form.Get("email"),
			Phone:     r.Form.Get("phone"),
			StartDate: startDate,
			EndDate:   endDate,
		}

		resvId, err := repo.DB.InsertReservation(reservation)

		if err != nil {
			http.Error(w, "Failed to Make Reservation.", http.StatusInternalServerError)
			log.Println("Failed to Write in Database", err)
			return
		}

		if err != nil {
			http.Error(w, "Failed to Make Reservation.", http.StatusInternalServerError)
			log.Println("Failed to Fectch ID from Database", err)
			return
		}

		roomRestriction := &models.RoomRestriction{
			RoomId:        1,
			ReservationId: int(resvId),
			RestrictionId: 1,
			StartDate:     startDate,
			EndDate:       endDate,
		}

		_, err = repo.DB.InsertRoomRestriction(roomRestriction)

		if err != nil {
			http.Error(w, "Failed to Make Reservation.", http.StatusInternalServerError)
			log.Println("Failed to Create Restriction to Database", err)
			_, err := repo.DB.DeleteReservation(resvId)
			if err != nil {
				log.Println("Could not Delete from Datavase", err)
			}
			return
		}

		repo.App.Session.Put(r.Context(), "reservation" /*+resvId.String()*/, r.Form)

		log.Println("POST:", r.RemoteAddr, "/book-now HTTP: 200")
		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
		return
	}

	//Django Style
	data := &struct {
		Id        string
		CSRFtoken string
	}{
		chi.URLParam(r, "id"),
		nosurf.Token(r),
	}

	render.Render(w, "make-reservation.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/book-now HTTP: 200")
}

func (repo *Repository) ReservationSumary(w http.ResponseWriter, r *http.Request) {
	values, ok := repo.App.Session.Get(r.Context(), "reservation").(url.Values)
	res := make(map[string]string)

	if !ok {
		fmt.Fprintf(w, "No Recent Reservations")
		log.Println("GET:", r.RemoteAddr, "/resarvation-summary HTTP: 200")
		return
	}

	for i, v := range values {
		if i == "csrf_token" {
			continue
		}

		res[i] = v[0]

	}

	repo.App.Session.Remove(r.Context(), "reservation")

	render.Render(w, "res-sum.page.tmpl", &struct{ DataMap map[string]string }{res})
	log.Println("GET:", r.RemoteAddr, "/resarvation-summary HTTP: 200")

}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	//Django Style
	data := &struct {
		Name      string
		SurName   string
		Age       int
		CSRFtoken string
	}{
		"Nick",
		"Andreou",
		29,
		nosurf.Token(r),
	}

	render.Render(w, "contact.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/contact HTTP: 200")
}
