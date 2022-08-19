package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aiym182/booking/internal/config"
	"github.com/aiym182/booking/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager

var testApp config.Config

func TestMain(m *testing.M) {

	//Register what am I going to put in the session. (things like struct)
	gob.Register(models.Reservation{})

	//change this to true in production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct {
}

func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWriter) Write(b []byte) (int, error) {

	length := len(b)
	return length, nil
}

func (mw *myWriter) WriteHeader(statuscode int) {

}
