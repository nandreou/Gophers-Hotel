package config

import (
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	Prd      bool
	Template map[string]*template.Template
	Session  *scs.SessionManager
}
