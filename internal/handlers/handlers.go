package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/driver"
	"github.com/Joshua-Russel/bookings/internal/forms"
	"github.com/Joshua-Russel/bookings/internal/helpers"
	"github.com/Joshua-Russel/bookings/internal/models"
	"github.com/Joshua-Russel/bookings/internal/render"
	"github.com/Joshua-Russel/bookings/internal/repository"
	"github.com/Joshua-Russel/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
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
func NewTestRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewTestingRepo(app),
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
		repo.App.Session.Put(r.Context(), "error", errors.New("reservation not found in session"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := repo.DB.GetRoomById(res.RoomID)

	if err != nil {
		repo.App.Session.Put(r.Context(), "error", errors.New("Room not found"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed
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
		repo.App.Session.Put(r.Context(), "error", "Cannot Parse Form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      models.Room{RoomName: r.Form.Get("room_name")},
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.IsValid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		stringMap := make(map[string]string)
		stringMap["start_date"] = startDate.Format(layout)
		stringMap["end_date"] = endDate.Format(layout)
		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	reservationId, err := repo.DB.InsertReservations(reservation)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Cannot insert reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: reservationId,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestrictions(restriction)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Cannot Insert RoomRestriction")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	htmlMsg := fmt.Sprintf(`
				<strong>Reservation Confirmation</strong><br>
				Dear %s,<br>
				This is to confirm your reservation from %s to %s  
`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	email := models.MailData{
		To:       reservation.Email,
		From:     "russ@163.com",
		Subject:  "Reservation confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	repo.App.MailChan <- email

	htmlMsg = fmt.Sprintf(`
				<strong>Reservation Notification</strong><br>
			A Reservation has been made for %s from %s to %s
`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	email = models.MailData{
		To:      "russ@163.com",
		From:    "russ@163.com",
		Subject: "Reservation Notification",
		Content: htmlMsg,
	}
	repo.App.MailChan <- email

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
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)

	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Cannot Parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Cannot Parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Cannot get availability")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	RoomID    string `json:"roomId"`
}

func (repo *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonData{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)

	endDate, _ := time.Parse(layout, ed)

	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := repo.DB.SearchAvailabilityByDatesAndRoomId(startDate, endDate, roomId)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonData{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	data := jsonData{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomId),
	}
	out, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Pop(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "cannot get item from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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
	//roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}

	exploded := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(exploded[2])
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.RoomID = roomId
	repo.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	var res models.Reservation

	startDate, _ := time.Parse("2006-01-02", sd)

	endDate, _ := time.Parse("2006-01-02", ed)

	room, err := repo.DB.GetRoomById(roomId)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.RoomID = roomId
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName
	repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (repo *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (repo *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.IsValid() {

		render.Template(w, r, "login.page.gohtml", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := repo.DB.Authenticate(email, password)
	if err != nil {

		repo.App.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(r.Context(), "user_id", id)
	repo.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/	", http.StatusSeeOther)
}

func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.Destroy(r.Context())
	_ = repo.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (repo *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})

}

func (repo *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-all-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})

}

func (repo *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	res, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "admin-reservation-show.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

func (repo *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")
	data := make(map[string]interface{})
	data["now"] = now
	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["last_month"] = lastMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := repo.DB.AllRooms()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)
		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := repo.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)

		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, y := range restrictions {
			if y.ReservationID > 0 {

				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {

					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID

			}
		}
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
		repo.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}
	render.Template(w, r, "admin-reservations-calendar.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})

}
func (repo *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = repo.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	month := r.Form.Get("month")
	year := r.Form.Get("year")

	repo.App.Session.Put(r.Context(), "flash", "Changes saved")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}
func (repo *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = repo.DB.UpdateProcessedForReservation(id, 1)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	repo.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)

	}
}

func (repo *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = repo.DB.DeleteReservation(id)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	repo.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

func (repo *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	rooms, err := repo.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	for _, x := range rooms {
		curMap := repo.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			if val, ok := curMap[name]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						err := repo.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			err := repo.DB.InsertBlockForRoom(roomId, t)
			if err != nil {
				log.Println(err)
			}
		}
	}
	repo.App.Session.Put(r.Context(), "flash", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
