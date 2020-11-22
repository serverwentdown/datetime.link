package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/serverwentdown/datetime.link/data"
	"go.uber.org/zap"
)

// ErrNoTemplate is returned when a template was not found on the server
var ErrNoTemplate = errors.New("missing template")

// Datetime is the main application server
type Datetime struct {
	*http.ServeMux
	tmpl   *template.Template
	cities map[string]*data.City
}

// NewDatetime creates an application instance. It assumes certain resources
// like templates and data exist
func NewDatetime() (*Datetime, error) {
	// Data
	tmpl, err := template.New("templates").Funcs(templateFuncs).ParseGlob("templates/*")
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
	mux.HandleFunc("/search", app.search)
	mux.HandleFunc("/", app.index)

	return app, nil
}

type appRequest struct {
	App Datetime
	Req Request
}

type appSearch struct {
	App    Datetime
	Search []*data.City
}

// index handles all incoming page requests
func (app Datetime) index(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmpl := app.loadTemplate("index", w, req)
	if tmpl == nil {
		return
	}
	if req.Method == http.MethodHead {
		return
	}

	request := Request{}
	if req.URL.Path != "/" {
		request, err = ParseRequest(req.URL)
		if err != nil {
			l.Debug("parse failed", zap.Error(err))
			app.error(HTTPError{http.StatusBadRequest, err}, w, req)
			return
		}
	}

	l.Debug("rendering template", zap.Reflect("request", request))
	err = tmpl.Execute(w, appRequest{app, request})
	if err != nil {
		l.Error("templating failed", zap.Error(err))
		app.templateError(HTTPError{http.StatusInternalServerError, err}, w, req)
		return
	}
}

// search handles zone search queries
func (app Datetime) search(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmpl := app.loadTemplate("search", w, req)
	if tmpl == nil {
		return
	}
	if req.Method == http.MethodHead {
		return
	}

	// TODO: do search
	query := req.URL.Query()
	search, err := FullSearchCities(app.cities, query.Get("zone"))
	if err != nil {
		l.Error("search failed", zap.Error(err))
		app.error(HTTPError{http.StatusInternalServerError, err}, w, req)
		return
	}

	l.Debug("rendering template", zap.Reflect("search", search))
	err = tmpl.Execute(w, appSearch{app, search})
	if err != nil {
		l.Error("templating failed", zap.Error(err))
		app.templateError(HTTPError{http.StatusInternalServerError, err}, w, req)
		return
	}
}
