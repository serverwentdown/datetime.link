package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"go.uber.org/zap"
)

// ErrTemplateNotFound is thrown when a template with the requested MIME type
// is not found.
var ErrTemplateNotFound = errors.New("unable to find template")

// loadTemplate returns a matching template for the request. It also causes an
// error if the template is not found or the Accept parameters are incorrect.
func (app Datetime) loadTemplate(name string, w http.ResponseWriter, req *http.Request) *template.Template {
	accept := req.Header.Get("Accept")
	tmpl, contentType, acceptable := app.chooseTemplate(accept, name)
	if !acceptable {
		if name == "error" {
			app.simpleError(HTTPError{http.StatusNotAcceptable, nil}, w, req)
			return nil
		}
		app.error(HTTPError{http.StatusNotAcceptable, nil}, w, req)
		return nil
	}
	if tmpl == nil {
		err := fmt.Errorf("%w \"%s\" for \"%s\"", ErrTemplateNotFound, name, accept)
		l.Warn("unable to find template", zap.Error(err), zap.String("name", name), zap.String("accept", accept))
		if name == "error" {
			app.simpleError(HTTPError{http.StatusNotAcceptable, err}, w, req)
		}
		app.error(HTTPError{http.StatusNotAcceptable, err}, w, req)
		//app.simpleError(HTTPError{http.StatusInternalServerError, ErrNoTemplate}, w, req)
		return nil
	}
	w.Header().Set("Content-Type", contentType)
	return tmpl
}

// chooseTemplate returns a template based on the accepted mime types from the
// client, and if a template cannot be found it returns a nil template.
func (app Datetime) chooseTemplate(accept string, name string) (t *template.Template, contentType string, acceptable bool) {
	acceptable = true
	switch chooseResponseType(accept) {
	case responsePlain:
		t = app.tmpl.Lookup(name + ".txt")
		contentType = "text/plain"
	case responseHTML:
		t = app.tmpl.Lookup(name + ".html")
		contentType = "text/html"
	case responseJSON:
		t = app.tmpl.Lookup(name + ".json")
		contentType = "application/json"
	case responseAny:
		t = app.tmpl.Lookup(name + ".txt")
		contentType = "text/plain"
	case responseUnknown:
		acceptable = false
		return
	}
	return
}
