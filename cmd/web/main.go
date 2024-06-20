package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/handlers"
	"github.com/Joshua-Russel/bookings/internal/helpers"
	"github.com/Joshua-Russel/bookings/internal/models"
	"github.com/Joshua-Russel/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNum = ":8000"

var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var app config.AppConfig

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Serving on port %s", portNum))

	serve := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	_ = serve.ListenAndServe()

}
func run() error {

	infoLog = log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	app.Session = session

	tempSet, err := render.CreateTemplate()
	if err != nil {
		log.Fatal("cannot create template:", err)
		return err
	}
	app.TmplCache = tempSet
	app.UseCache = false
	app.InProduction = false

	render.NewTemplate(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.Newhelpers(&app)
	return nil
}
