package main

import (
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/alexedwards/scs/v2"
	"nicksrepo.com/nick/pkg/config"
	"nicksrepo.com/nick/pkg/database"
	handlerfuntions "nicksrepo.com/nick/pkg/handlerFuntions"
	"nicksrepo.com/nick/pkg/middlewares"
	"nicksrepo.com/nick/pkg/models"
	"nicksrepo.com/nick/pkg/render"
	"nicksrepo.com/nick/pkg/routers"
)

const (
	dbData = "host=localhost port=5432 dbname=mydb user=nick password=password"
)

func main() {

	var app config.App

	db, srv, err := run(&app)

	if err != nil {
		log.Fatal("Failed to Run Server check Template Caching")
	}
	defer db.SQL.Close()

	log.Println("Server is UP Listening to :8000")

	if err := srv.ListenAndServe(); err != nil {
		log.Println("Unfortunately Could not Start server")
	}

}

func run(app *config.App) (*database.DB, *http.Server, error) {
	var cacheError error

	gob.Register(url.Values{})
	gob.Register(models.Reservation{})

	//Pass App config In Pakages
	render.NewCache(app)
	middlewares.PassConfigToMidPkg(app)

	//Cache the Templates in Servers Memory
	app.Template, cacheError = render.CreateTmplCache()

	if cacheError != nil {
		return nil, nil, errors.New("sdf") //cacheError
	}

	app.Prd = false

	//Create Session Manager and session configs
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Prd

	app.Session = session

	//Connect to DataBase
	db, err := database.ConnectToDatabase(dbData)
	if err != nil {
		log.Fatal("Could not connect to Database: ", err)
	}

	log.Println("Checking Database Connection...")

	err = db.SQL.Ping()
	if err != nil {
		log.Fatal("Could not connect to Database: ", err)
	} else {
		log.Println("Connected to Database")
	}

	//Create Repository
	handlerfuntions.NewRepository(app, db)

	return db, &http.Server{
		Addr:    ":8000",
		Handler: routers.Routers(),
	}, nil
}
