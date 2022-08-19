package main

import (
	"testing"

	"github.com/aiym182/booking/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {

	app := &config.Config{}

	mux := Routes(app)

	switch v := mux.(type) {

	case *chi.Mux:

	default:
		t.Errorf("Type is not *chi.Mux, instead it is %T", v)
	}
}
