package handlers

import (
	"fmt"
	"net/http"

	"github.com/aiym182/booking/pkg/config"
	"github.com/aiym182/booking/pkg/models"
	"github.com/aiym182/booking/pkg/render"
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
	render.RenderTemplate(rw, "home.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(rw, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
