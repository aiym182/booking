package render

import (
	"net/http"
	"testing"

	"github.com/aiym182/booking/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Errorf("Flash value of 123 not found in session")
	}
}

func TestTemplateCache(t *testing.T) {
	tc, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}
	var rw myWriter

	err = Template(&rw, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}
	err = Template(&rw, r, "non-existant.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some", nil)

	if err != nil {
		return nil, err
	}
	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
