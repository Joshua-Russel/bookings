package main

import (
	"fmt"
	"github.com/Joshua-Russel/bookings/pkg/config"
	"github.com/Joshua-Russel/bookings/pkg/handlers"
	"github.com/Joshua-Russel/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNum = ":8000"

var session *scs.SessionManager

var app config.AppConfig

func main() {

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	app.Session = session

	tempSet, err := render.CreateTemplate()
	if err != nil {
		log.Fatal("cannot create template:", err)
	}
	app.TmplCache = tempSet
	app.UseCache = false
	app.InProduction = false

	render.NewTemplate(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println("serving on port", portNum)

	serve := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	_ = serve.ListenAndServe()

}
