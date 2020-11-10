package main

import (
	"html/template"
	"net/http"

	"go.uber.org/zap"
)

var templateFuncs = map[string]interface{}{
	"statusText": templateFuncStatusText,
	"thisIsSafe": templateFuncThisIsSafe,
	// Formatting
	"formatOffset": templateFuncFormatOffset,
	// Logic
	"resolveZone": templateFuncResolveZone,
}

func templateFuncStatusText(s int) string {
	return http.StatusText(s)
}
func templateFuncThisIsSafe(s string) template.HTML {
	return template.HTML(s)
}

func templateFuncFormatOffset(offset int) string {
	return FormatZoneOffset(offset)
}

// ResolvedZone holds a resolved zone or an error
type ResolvedZone struct {
	Zone
	Error error
}

func templateFuncResolveZone(app Datetime, zone string) ResolvedZone {
	z, err := ResolveZone(app.cities, zone)
	if err != nil {
		l.Debug("unable to resolve zone", zap.Reflect("zone", zone), zap.Error(err))
	}
	return ResolvedZone{z, err}
}
