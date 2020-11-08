package main

import (
	"html/template"
	"net/http"

	"github.com/serverwentdown/datetime.link/data"
	"go.uber.org/zap"
)

// Datetime is the main application server
type Datetime struct {
	*http.ServeMux
	tmpl   *template.Template
	cities map[string]*data.City
}

// NewDatetime creates an application instance. It assumes certain resources
// like templates and data exist.
func NewDatetime() (*Datetime, error) {
	// Data
	tmpl, err := template.ParseGlob("templates/*")
	if err != nil {
		return nil, err
	}
	cities, err := data.ReadCities()
	if err != nil {
		return nil, err
	}

	// Mux
	mux := http.NewServeMux()
	app := &Datetime{mux, tmpl, cities}

	// Routes
	mux.Handle("/data/", http.FileServer(http.Dir(".")))
	mux.Handle("/js/", http.FileServer(http.Dir("assets")))
	mux.Handle("/css/", http.FileServer(http.Dir("assets")))
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("assets")))
	mux.HandleFunc("/", app.index)

	return app, nil
}

// index handles all incoming page requests
func (app Datetime) index(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accept := req.Header.Get("Accept")
	tmpl, acceptable := app.chooseTemplate(accept)
	if !acceptable {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if tmpl == nil {
		l.Error("unable to find template", zap.String("accept", accept))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodHead {
		return
	}

	l.Debug("", zap.Reflect("url", req.URL))
	err = tmpl.Execute(w, nil)
	if err != nil {
		l.Error("templating failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError) // Usually this will fail
		return
	}
}

// chooseTemplate returns a template based on the accepted mime types from the
// client, and if a template cannot be found it returns a nil template.
func (app Datetime) chooseTemplate(accept string) (t *template.Template, acceptable bool) {
	responseType := chooseResponseType(accept)
	templateName := ""
	switch responseType {
	case responsePlain:
		templateName = "index.txt"
	case responseHTML:
		templateName = "index.html"
	case responseAny:
		templateName = "index.txt"
	case responseUnknown:
		return nil, false
	}
	return app.tmpl.Lookup(templateName), true
}
