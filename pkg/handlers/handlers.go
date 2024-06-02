package handlers

import (
	"github.com/Joshua-Russel/bookings/pkg/config"
	"github.com/Joshua-Russel/bookings/pkg/models"
	"github.com/Joshua-Russel/bookings/pkg/render"
	"log"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", addr)
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{
		Flash: "sucdess",
	})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	td := make(map[string]string)
	td["test"] = "hello,again"
	remoteIp := repo.App.Session.GetString(r.Context(), "remote_ip")
	td["remote_ip"] = remoteIp
	log.Println(remoteIp)
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: td,
	})

}
