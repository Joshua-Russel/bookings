package main

import (
	"github.com/Joshua-Russel/bookings/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app.InProduction,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}
func SessionLoad(next http.Handler) http.Handler {

	return session.LoadAndSave(next)
}
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Login to continue")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return

		}
		next.ServeHTTP(w, r)
	})
}
