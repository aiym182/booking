package models

import "github.com/aiym182/booking/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]any
	CSRFToken string //cross site request forgery token
	Flash     string
	Warning   string
	Error     string
	Forms     *forms.Form
}
