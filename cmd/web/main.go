package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/handlers"
	"github.com/Joshua-Russel/bookings/internal/models"
	"github.com/Joshua-Russel/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNum = ":8000"

var session *scs.SessionManager

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
	return nil
}
