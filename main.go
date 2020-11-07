package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
)

var listen string
var tmpl *template.Template

func main() {
	var err error

	flag.StringVar(&listen, "listen", ":8000", "Listen address")
	flag.Parse()

	server := &http.Server{
		Addr:         listen,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	tmpl, err = template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}

	http.Handle("/data/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir("assets")))
	http.Handle("/css/", http.FileServer(http.Dir("assets")))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/", index)

	log.Printf("Listening on %v", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accept := req.Header.Get("Accept")
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
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	t := tmpl.Lookup(templateName)
	if t == nil {
		log.Printf("Unable to find index template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Method == http.MethodHead {
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Error: %v", err)
		// Usually, the following will fail
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
