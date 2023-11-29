package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/driver"
	"github.com/aiym182/booking/internal/forms"
	"github.com/aiym182/booking/internal/models"
	"github.com/aiym182/booking/internal/render"
	"github.com/aiym182/booking/internal/repository"
	"github.com/aiym182/booking/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

//Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.Config
	DB  repository.DatabaseRepo
}

// TemplateData holds data sent from handlers to templates

// Create a new Repository
func NewRepo(a *config.Config, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// only use this for testing purposes
func NewTestRepo(a *config.Config) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// Sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	// saving remote ip adress from user to session
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.Template(rw, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {

	// perform some logic
	stringMap := make(map[string]string)
	//getting remote ip from session and save it to stringMap
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	stringMap["test"] = "Hello world again"

	fmt.Println(stringMap["remote_ip"])
	// send the data to the template
	render.Template(rw, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders reservation page and display form.
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = res

	render.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Forms:     forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of reservation form
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form !")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLen("first_name", 3)
	form.MinLen("last_name", 3)
	form.ValidateEmail("email")

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	if !form.Valid() {

		data := make(map[string]any)
		data["reservation"] = reservation
		// http.Error(rw, "my own error message", http.StatusSeeOther)
		render.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
			Forms:     form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return

	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room_restriction into database")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return

	}

	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s, <br>
	This is confirm your reservaiton from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-2"))
	// send notification -first to guest

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	htmlMessage = fmt.Sprintf(`
	<strong>Reservation Notification</strong><br>
	A reservation has been made for %s from %s to %s.
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.MailData{
		To:      "me@here.com",
		From:    "me@here.com",
		Subject: "Reservation Notification",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	//Saving info from above(form) through session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}

//Generals renders the room page

func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "generals.page.tmpl", &models.TemplateData{})
}

//Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "majors.page.tmpl", &models.TemplateData{})
}

//Majors renders the room page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "search-available.page.tmpl", &models.TemplateData{})
}

//Post availability posts something cool
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabiltyForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find search availability for all rooms")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "There are no available rooms in these dates")
		http.Redirect(rw, r, "search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]any)
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(rw, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJson handles request for availability and send JSON response

func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		// can't parse form so return appropriate json

		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "\t")
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	if err != nil {
		// can't parse form so return appropriate json

		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting database",
		}
		out, _ := json.MarshalIndent(resp, "", "\t")
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "\t")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

//Majors renders the room page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummer displays reservation summary
func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	//models.Reservation is type Assert
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = reservation

	render.Template(rw, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})

}

// Choose room displays available rooms
func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {

	// roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	splitted := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(splitted[2])

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse roomID from url parameter")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes url parameters, builds a sessional variable , and takes user to make reservation screen
func (m *Repository) BookRoom(rw http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "login.page.tmpl", &models.TemplateData{
		Forms: forms.New(nil),
	})
}

//postshowlogin handles the user in
func (m *Repository) PostShowLogin(rw http.ResponseWriter, r *http.Request) {

	// renewtoken updates session data to have new session data and reset session timeout
	// it is always to good to renew token when log in
	m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")

		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.ValidateEmail("email")
	if !form.Valid() {
		render.Template(rw, r, "login.page.tmpl", &models.TemplateData{
			Forms: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(rw, r, "/", http.StatusSeeOther)

}

// Logout logs user out
func (m *Repository) Logout(rw http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
}

//AdminDashBoard shows dashboard in all admin tool
func (m *Repository) AdminDashBoard(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

//AdminNewReservations shows all new reservations in all admin tool
func (m *Repository) AdminNewReservations(rw http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error getting AllNewReservations from Database")
		return
	}
	data := make(map[string]any)
	data["reservations"] = reservations
	render.Template(rw, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

//AdminAllReservations shows all reservations in admin tool
func (m *Repository) AdminAllReservations(rw http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error getting AllReservations from Database")
		return
	}
	data := make(map[string]any)
	data["reservations"] = reservations
	render.Template(rw, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

//AdminShowReservation shows the reservation in admin tool
func (m *Repository) AdminShowReservation(rw http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/") // r.RequestURI gives the url by client

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error converting url to int")
		return
	}

	src := exploded[3]

	StringMap := make(map[string]string)
	StringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	StringMap["month"] = month
	StringMap["year"] = year

	//Get reservation from database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error getting data from database")
		return
	}

	data := make(map[string]any)

	data["reservation"] = res

	render.Template(rw, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
		Data:      data,
		Forms:     forms.New(nil),
		StringMap: StringMap,
	})

}

//AdminShowReservation shows the reservatioin in admin tool
func (m *Repository) AdminPostShowReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't parse form")
		return
	}

	exploded := strings.Split(r.RequestURI, "/") // r.RequestURI gives the url by client

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Got invalid Url")
		return
	}
	src := exploded[3]
	StringMap := make(map[string]string)
	StringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Got invalid id")
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Got invalid reservation")
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "change saved")

	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//AdminReservationCalendar shows reservation calendar in admin tool
func (m *Repository) AdminReservationCalendar(rw http.ResponseWriter, r *http.Request) {

	// assume that there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]any)
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and the last days of the month

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error getting AllRooms from database")
		return
	}
	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps

		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}
		// get all the restrictions for the current room

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)

		if err != nil {
			m.App.Session.Put(r.Context(), "error", "error calling GetRestrictionsForRoomByDate from database")
			return
		}
		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservation
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {

				//it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
		// need to save block map in session in order to compare blocked data after the change
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(rw, r, "admin-reservation-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})

}

//AdminProcessReservation marks a reservation as processed
func (m *Repository) AdminProcessReservation(rw http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	err := m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		log.Println(err)
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)

	}

}

//AdminDeleteReservation deletes reservation
func (m *Repository) AdminDeleteReservation(rw http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// AdminPostReservationCalendar handles post of reservation calendar
func (m *Repository) AdminPostReservationCalendar(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "unable to parse form")
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process blocks

	rooms, err := m.DB.AllRooms()

	if err != nil {
		log.Println("Error getting AllRooms database")
		m.App.Session.Put(r.Context(), "error", "Error getting AllRooms database")
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		// Get the block map from the session. Loop through entire map, if we have an entry in the map
		// that does not exist in our posted data, and if the restriciton id > 0, then it is a block we need to remove
		// remove.

		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)

		for name, value := range curMap {

			// ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to values >0, and that are not in the form post
				// the rest are just placeholders for days without blocks

				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name), r) {
						// delete the restriction by id

						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println("Error getting DeleteBlockByID database")

							m.App.Session.Put(r.Context(), "error", "Error getting DeleteBlockByID database")
							return
						}
					}
				}

			}

		}
	}

	// now handle new blocks
	for name := range r.PostForm {

		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert a new block
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println("Error getting InsertBlockForRoom database")
				m.App.Session.Put(r.Context(), "error", "Error getting InsertBlockForRoom database")
				return
			}

		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(rw, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
