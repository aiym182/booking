package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/driver"
	"github.com/aiym182/booking/internal/handlers"
	"github.com/aiym182/booking/internal/helpers"
	"github.com/aiym182/booking/internal/models"
	"github.com/aiym182/booking/internal/render"
	"github.com/alexedwards/scs/v2"
)

const webport = ":8080"

var app = &config.Config{}
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	// starting mail listener...
	log.Println("Starting mail listener...")
	listenForMail()

	log.Printf("Web application on port %s\n", webport)

	serve := &http.Server{
		Addr:    webport,
		Handler: Routes(app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {

	//Register what am I going to put in the session. (things like struct)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change this to true in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	//connect to DB

	log.Println("Connecting to database..")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=Bookings user=brandonlee password=")
	if err != nil {
		log.Fatalf("Cannot connect to database : %v", err)
	}

	log.Println("Successfully connected to database!!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(app, db)

	handlers.NewHandlers(repo)
	render.NewRenderer(app)
	helpers.NewHelpers(app)

	return db, nil
}
