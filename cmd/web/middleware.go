package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(rw, r)
	})
}

//No Surf adds csrf protections to all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, //In production, you have to switch to yes
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler

}

// loads and save the session on every request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
