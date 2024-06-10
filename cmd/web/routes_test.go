package main

import (
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//Nothing
	default:
		t.Error(fmt.Sprintf("Type did not match chi.Mux ,%T", v))
	}
}
