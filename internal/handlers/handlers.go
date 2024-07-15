package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/driver"
	"github.com/Joshua-Russel/bookings/internal/forms"
	"github.com/Joshua-Russel/bookings/internal/helpers"
	"github.com/Joshua-Russel/bookings/internal/models"
	"github.com/Joshua-Russel/bookings/internal/render"
	"github.com/Joshua-Russel/bookings/internal/repository"
	"github.com/Joshua-Russel/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})

}
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gohtml", &models.TemplateData{})

}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gohtml", &models.TemplateData{})
}
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Reservation not found in session"))
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	room, err := repo.DB.GetRoomById(res.RoomID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	repo.App.Session.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringmap,
	})
}
func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Reservation not found in session"))
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	if !form.IsValid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	reservationId, err := repo.DB.InsertReservations(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservationId,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestrictions(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	repo.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	start := r.PostFormValue("start_date")
	end := r.PostFormValue("end_date")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		repo.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/srch-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	repo.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

type jsonData struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := repo.DB.SearchAvailabilityByDatesAndRoomId(startDate, endDate, roomId)

	data := jsonData{
		OK:      available,
		Message: "",
	}
	js, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Pop(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.ErrorLog.Println("cannot get item from session")
		repo.App.Session.Put(r.Context(), "error", "cannot get item from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringmap,
	})
}

func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}
	res.RoomID = roomId
	repo.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
