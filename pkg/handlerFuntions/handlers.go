package handlerfuntions

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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

	render.Render(w, "home.page.tmpl", nil)
	log.Println("GET:", r.RemoteAddr, "/ HTTP: 200")
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

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

	if r.Method == "POST" {

		resData := &models.Reservation{}

		//PARSING FORM VALUES
		err := r.ParseForm()

		if err != nil {
			log.Println(err)
			return
		}

		//Convertingg date to time.Time
		dateFormat := "2006-01-02"

		start, err := time.Parse(dateFormat, r.Form.Get("start"))

		if err != nil {
			http.Error(w, "Failed to Make Reservation Please use 2024-01-02 StartDate Format", http.StatusBadRequest)
			log.Println("Failed to Write in StartDate Database", start, "\n", err)
			return
		}

		end, err := time.Parse(dateFormat, r.Form.Get("end"))

		if err != nil {
			http.Error(w, "Failed to Make Reservation Please use 2024-01-02 EndDate Format", http.StatusBadRequest)
			log.Println("Failed to Write in EndDate Database", end, "\n", err)
			return
		}

		//Adding dates to Rservation model
		resData.StartDate = start
		resData.EndDate = end

		rooms, err := repo.DB.AvailabilitySearch(r.Form.Get("start"), r.Form.Get("end"))

		if err != nil {
			log.Println(err)
		}

		resId := uuid.New().String()

		data := &struct {
			Rooms     []*models.SearchAvailabilityModel
			ResId     string
			CSRFtoken string
		}{
			rooms,
			resId,
			nosurf.Token(r),
		}

		//Creating Id for Reservation Request

		repo.App.Session.Put(r.Context(), resId, resData)

		log.Println("POST:", r.RemoteAddr, "/search-availability HTTP: 200")
		render.Render(w, "availability.page.tmpl", data)

		return
	}

	data := &struct {
		CSRFtoken string
	}{
		nosurf.Token(r),
	}

	render.Render(w, "search-availability.page.tmpl", data)
	log.Println("GET:", r.RemoteAddr, "/search-availability HTTP: 200")
}

func (repo *Repository) BookNow(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		query := r.URL.Query().Get("ri")
		fmt.Fprintf(w, "%s  ", query)
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error on Parshing Form data Try Again Later %s", err)
			return
		}

		form := forms.NewForm(r.Form)

		//SimpleValidateTest Form Validation
		for field := range r.Form {
			form.Has(field)
		}

		res, ok := repo.App.Session.Get(r.Context(), query).(models.Reservation)

		if !ok {
			log.Println("Something Went Wrong Session Data Is Empty")
			fmt.Fprintf(w, "Something Went Wrong Session Data Is Empty")
			return
		}

		//Last check of availability
		id := chi.URLParam(r, "id")

		form.ValidForm, _ = repo.DB.LastAvailabilitySearch(id, res.StartDate, res.EndDate)

		//Check Form Validation
		if !form.IsValid() {
			log.Println("POST:", r.RemoteAddr, "/book-now HTTP: 400")
			http.Redirect(w, r, "/book-now/"+id+"?ri="+query, http.StatusSeeOther)
			return
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			log.Println("Something Went Wrong")
			fmt.Fprintf(w, "Something went wrong")
			return
		}

		reservation := &models.Reservation{
			RoomId:    idInt,
			FirstName: r.Form.Get("first_name"),
			LastName:  r.Form.Get("last_name"),
			Email:     r.Form.Get("email"),
			Phone:     r.Form.Get("phone"),
			StartDate: res.StartDate,
			EndDate:   res.EndDate,
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
			RoomId:        idInt,
			ReservationId: int(resvId),
			RestrictionId: 1,
			StartDate:     res.StartDate,
			EndDate:       res.EndDate,
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

		//Updating session data
		repo.App.Session.Put(r.Context(), query, reservation)

		log.Println("POST:", r.RemoteAddr, "/book-now HTTP: 200")
		http.Redirect(w, r, "/reservation-summary?ri="+query, http.StatusSeeOther)
		return
	}

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
	query := r.URL.Query().Get("ri")
	res, ok := repo.App.Session.Get(r.Context(), query).(models.Reservation)

	if !ok {
		fmt.Fprintf(w, "No Recent Reservations")
		log.Println("GET:", r.RemoteAddr, "/resarvation-summary HTTP: 200")
		return
	}

	fmt.Println(res)

	repo.App.Session.Remove(r.Context(), query)

	render.Render(w, "res-sum.page.tmpl", res)
	log.Println("GET:", r.RemoteAddr, "/resarvation-summary HTTP: 200")

}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

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
