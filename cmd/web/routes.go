package main

import (
	"github.com/Joshua-Russel/bookings/pkg/config"
	"github.com/Joshua-Russel/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/rooms/majors", handlers.Repo.Majors)
	mux.Get("/rooms/generals", handlers.Repo.Generals)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/srch-availability", handlers.Repo.SearchAvailability)
	mux.Post("/srch-availability", handlers.Repo.Availability)
	mux.Post("/srch-availability-json", handlers.Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
