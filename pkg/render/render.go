package render

import (
	"errors"
	"net/http"
	"path/filepath"
	"text/template"

	"nicksrepo.com/nick/pkg/config"
)

var app *config.App
var functions = template.FuncMap{}
var pathToTemplates = "../templates/*"
var pathToTestTemplates = "../../templates/*"

func NewCache(a *config.App) {
	app = a
}

func CreateTmplCache() (map[string]*template.Template, error) {
	Cache := map[string]*template.Template{}
	pages, err := filepath.Glob(pathToTemplates)

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		newTmpl := template.New(name)
		newCache, err := newTmpl.Funcs(functions).ParseFiles(page)

		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob("../templates/*.layout.tmpl")
		if err != nil {
			return Cache, err
		}

		if len(matches) > 0 {
			newCache, err = newCache.ParseGlob("../templates/*.layout.tmpl")
			if err != nil {
				return Cache, err
			}
		}
		Cache[name] = newCache

	}
	return Cache, nil
}

func Render(w http.ResponseWriter, pg string, data interface{}) error {
	if !app.Prd {

		tmpl, err := CreateTmplCache()
		if err != nil {
			http.Error(w, "Sorry an Error accured", 500)
			return err
		}

		tmpl[pg].Execute(w, data)

	} else {
		_, ok := app.Template[pg]
		if !ok {
			http.Error(w, "Page Not Found", 404)
			return errors.New("error on Reading From Cache, template may not exist")
		} else {
			app.Template[pg].Execute(w, data)
		}

	}
	return nil
}
