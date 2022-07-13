package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// holds app config
type Config struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}
