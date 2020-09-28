package main

import (
	"encoding/json"
	"strings"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

const timezoneXML = "https://raw.githubusercontent.com/unicode-org/cldr/master/common/bcp47/timezone.xml"

func main() {
	resp, err := http.Get(timezoneXML)
	if err != nil {
		log.Fatalf("Failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to fetch: %v", err)
	}

	// Parse XML file
	data := &ldmlData{}
	err = xml.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}

	// Remap into different format
	timezonesData := make([]timezoneData, len(data.Keys))
	for i, t := range data.Keys {
		timezonesData[i] = timezoneData{
			Name: t.Name,
			Description: t.Description,
			Aliases: strings.Split(t.Alias, " "),
		}
	}

	// Encode JSON file
	b, err := json.Marshal(timezonesData)
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}

	// Write JSON file
	err = ioutil.WriteFile("js/bcp47timezone.json", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}
}

type ldmlData struct {
	Keys []ldmlType `xml:"keyword>key>type"`
}

type ldmlType struct {
	Name string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Alias string `xml:"alias,attr"` // NOTE: space-separated values
}

type timezoneData struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Aliases []string `json:"aliases"`
}
