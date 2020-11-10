package main

import (
	"net/http"

	"go.uber.org/zap"
)

// HTTPError is an error that can be rendered into a HTML page to follow HTTP
// status semantics
type HTTPError struct {
	Status int
	Error  error
}

func (app Datetime) error(httpErr HTTPError, w http.ResponseWriter, req *http.Request) {
	tmpl := app.loadTemplate("error", w, req)
	if tmpl == nil {
		return
	}

	w.WriteHeader(httpErr.Status)
	if req.Method == http.MethodHead {
		return
	}

	err := tmpl.Execute(w, httpErr)
	if err != nil {
		l.Error("template failed", zap.Error(err))
		app.templateError(HTTPError{http.StatusInternalServerError, err}, w, req)
		return
	}
}

func (app Datetime) simpleError(httpErr HTTPError, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(httpErr.Status)
	if req.Method == http.MethodHead {
		return
	}

	if httpErr.Error != nil {
		w.Write([]byte(http.StatusText(httpErr.Status) + ": " + httpErr.Error.Error()))
	} else {
		w.Write([]byte(http.StatusText(httpErr.Status)))
	}
}

func (app Datetime) templateError(httpErr HTTPError, w http.ResponseWriter, req *http.Request) {
	// Sadly, we probably already sent out the header
	//w.Header().Set("Content-Type", "text/plain")
	//w.WriteHeader(httpErr.Status)

	if httpErr.Error != nil {
		w.Write([]byte("\n" + http.StatusText(httpErr.Status) + ": " + httpErr.Error.Error()))
	} else {
		w.Write([]byte("\n" + http.StatusText(httpErr.Status)))
	}
}
