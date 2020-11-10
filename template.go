package main

import (
	"html/template"
	"net/http"

	"go.uber.org/zap"
)

// loadTemplate returns a matching template for the request. It also causes an
// error if the template is not found or the Accept parameters are incorrect.
func (app Datetime) loadTemplate(name string, w http.ResponseWriter, req *http.Request) *template.Template {
	accept := req.Header.Get("Accept")
	tmpl, contentType, acceptable := app.chooseTemplate(accept, name)
	if !acceptable {
		app.simpleError(HTTPError{http.StatusNotAcceptable, nil}, w, req)
		return nil
	}
	if tmpl == nil {
		l.Error("unable to find template", zap.String("name", name), zap.String("accept", accept))
		app.simpleError(HTTPError{http.StatusInternalServerError, ErrNoTemplate}, w, req)
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
	case responseAny:
		t = app.tmpl.Lookup(name + ".txt")
		contentType = "text/plain"
	case responseUnknown:
		acceptable = false
		return
	}
	return
}
