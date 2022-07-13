package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/forms"
	"github.com/aiym182/booking/internal/models"
	"github.com/aiym182/booking/internal/render"
)

//Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.Config
}

// TemplateData holds data sent from handlers to templates

// Create a new Repository
func NewRepo(a *config.Config) *Repository {
	return &Repository{App: a}
}

// Sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	// saving remote ip adress from user to session
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.RenderTemplate(rw, r, "home.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(rw, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders reservation page and display form.
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservations
	data := make(map[string]any)
	data["reservation"] = emptyReservation

	render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Forms: forms.New(nil),
		Data:  data,
	})
}

// PostReservation handles the posting of reservation form
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservations{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLen("first_name", 3)
	form.MinLen("last_name", 3)
	form.ValidatEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation

		render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
			Forms: form,
			Data:  data,
		})
		return
	}
	//Saving info from above(form) through session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}

//Generals renders the room page

func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "generals.page.tmpl", &models.TemplateData{})
}

//Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "majors.page.tmpl", &models.TemplateData{})
}

//Majors renders the room page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "search-available.page.tmpl", &models.TemplateData{})
}

//Post availability posts something cool
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	rw.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson handles request for availability and send JSON response

func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      false,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "\t")

	if err != nil {
		log.Println(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

//Majors renders the room page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	//models.Reservation is type Assert
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservations)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	}
	data := make(map[string]any)
	data["reservation"] = reservation

	render.RenderTemplate(rw, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})

}
