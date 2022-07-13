package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.Config

func NewTemplates(a *config.Config) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	//PopString() fetches string value for a given key and then delete it form session data.
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td

}

// Rendertemplate renders templates using html/templates
func RenderTemplate(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = RenderTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Cannot find any templates")
	}

	td = AddDefaultData(td, r)

	err = t.Execute(rw, td)

	if err != nil {
		fmt.Println(err)
	}

}

// create a template cache as a map
func RenderTemplateCache() (map[string]*template.Template, error) {

	fs, _ := filepath.Abs("../../template")

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", fs))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", fs))

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err := ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", fs))

			if err != nil {
				return myCache, nil
			}
			myCache[name] = ts
		}

	}

	return myCache, nil
}
