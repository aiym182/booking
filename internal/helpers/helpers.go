package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/aiym182/booking/internal/config"
)

var app *config.Config

// New Helpers sets up config for helpers
func NewHelpers(a *config.Config) {
	app = a
}

func ClientError(rw http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(rw, http.StatusText(status), status)
}
func ServerError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s \n %s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
