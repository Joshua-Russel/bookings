package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

type AppConfig struct {
	UseCache     bool
	TmplCache    map[string]*template.Template
	InProduction bool
	InfoLog      *log.Logger
	Session      *scs.SessionManager
}
