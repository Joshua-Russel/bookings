package main

import (
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/handlers"
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
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Get("/srch-availability", handlers.Repo.SearchAvailability)
	mux.Post("/srch-availability", handlers.Repo.Availability)
	mux.Post("/srch-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.Route("/admin", func(r chi.Router) {
		//r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)
		r.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		r.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		r.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		r.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)
		r.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		r.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		r.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		r.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
	})

	return mux
}
