package main

import (
	"fmt"
	"html/template"
	"testing"
)

func TestChooseTemplate(t *testing.T) {
	tmpl, err := template.New("templates").Funcs(templateFuncs).ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}
	app := &Datetime{tmpl: tmpl}

	type chooseTest struct {
		accept      string
		acceptable  bool
		contentType string
		template    string
	}
	tests := []chooseTest{
		{"text/html", true, "text/html", "index.html"},
		{"text/html;q=0.9,text/plain", true, "text/plain", "index.txt"},
		{"image/png", false, "", ""},
		{"*/*", true, "text/plain", "index.txt"},
	}

	for _, test := range tests {
		tmpl, contentType, acceptable := app.chooseTemplate(test.accept, "index")
		fn := fmt.Sprintf("chooseTemplate(\"%s\")", test.accept)
		if contentType != test.contentType {
			t.Errorf("%s; contentType = %v; wanted %v", fn, contentType, test.contentType)
		}
		if acceptable != test.acceptable {
			t.Errorf("%s; acceptable = %v; wanted %v", fn, acceptable, test.acceptable)
		}
		if tmpl != app.tmpl.Lookup(test.template) {
			t.Errorf("%s; tmpl = %v; wanted template for %v", fn, tmpl.Name(), test.template)
		}
	}
}
