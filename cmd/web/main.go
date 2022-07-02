package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aiym182/booking/pkg/config"
	"github.com/aiym182/booking/pkg/handlers"
	"github.com/aiym182/booking/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const webport = ":8080"

var app = &config.Config{}
var session *scs.SessionManager

func main() {

	//change this to true in production
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.RenderTemplateCache()
	if err != nil {
		log.Panic(err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(app)

	handlers.NewHandlers(repo)
	render.NewTemplates(app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	log.Printf("Web application on port %s\n", webport)

	// err = http.ListenAndServe(webport, nil)
	// if err != nil {
	// 	log.Panicln("Internal Server error : ", err)
	// }

	serve := &http.Server{
		Addr:    webport,
		Handler: Routes(app),
	}

	err = serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
