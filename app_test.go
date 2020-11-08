package main

import (
	"fmt"
	"html/template"
	"testing"
)

func TestChooseTemplate(t *testing.T) {
	tmpl, err := template.ParseGlob("templates/*")
	if err != nil {
		t.Errorf("unable to load templates: %v", err)
	}
	app := &Datetime{tmpl: tmpl}

	type chooseTest struct {
		accept     string
		acceptable bool
		template   string
	}
	tests := []chooseTest{
		{"text/html", true, "index.html"},
		{"text/html;q=0.9,text/plain", true, "index.txt"},
		{"image/png", false, ""},
		{"*/*", true, "index.txt"},
	}

	for _, test := range tests {
		tmpl, acceptable := app.chooseTemplate(test.accept)
		fn := fmt.Sprintf("chooseTemplate(\"%s\")", test.accept)
		if acceptable != test.acceptable {
			t.Errorf("%s; acceptable = %v; wanted %v", fn, acceptable, test.acceptable)
		}
		if tmpl != app.tmpl.Lookup(test.template) {
			t.Errorf("%s; tmpl = %v; wanted template for %v", fn, tmpl.Name(), test.template)
		}
	}
}
