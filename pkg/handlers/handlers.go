package handlers

import (
	"fmt"
	"github.com/Joshua-Russel/bookings/pkg/config"
	"github.com/Joshua-Russel/bookings/pkg/models"
	"github.com/Joshua-Russel/bookings/pkg/render"
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
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{})

}
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})

}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	cin := r.PostFormValue("start_date")

	cout := r.PostFormValue("end_date")

	n, _ := w.Write([]byte(fmt.Sprintf("the check in date is %s and check out date is %s", cin, cout)))
	fmt.Println(n)
}
