package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// holds app config
type Config struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
}
