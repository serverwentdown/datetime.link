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
		WriteTimeout: 10 * time.Second,
	}

	tmpl, err = template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}

	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", index)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accept := req.Header.Get("Accept")
	responseType := chooseResponseType(accept)

	switch responseType {
	case responsePlain:
		w.WriteHeader(http.StatusNotImplemented)
	case responseHTML:
		indexTmpl := tmpl.Lookup("index.html")
		if indexTmpl == nil {
			log.Printf("Unable to find index template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = indexTmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
