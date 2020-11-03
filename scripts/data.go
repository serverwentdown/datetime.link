/*
This script reads in GeoNames data and creates a table of IDs to city names and
timezones. The IDs are created from the ASCII city name, with administrative
division level 1 name and country code as disambiguation. Examples of IDs are:

	Singapore
	Ashland_Oregon_US
	Ashland_Mississippi_US

*/
package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// City represents a city that belongs inside an administrative division level 1
// and a country
type City struct {
	// Ref is the ASCII name of the city
	Ref string `json:"r"`
	// Name is the full UTF-8 name of the city
	Name           string   `json:"n"`
	AlternateNames []string `json:"an"`

	// Admin1Ref is the ASCII name of the administrative division level 1
	Admin1Ref  string `json:"a1r"`
	Admin1Name string `json:"a1n"`

	// CountryRef is the ISO-3166 2-letter country code
	CountryRef  string `json:"cr"`
	CountryName string `json:"cn"`
}

func main() {
}

/*

func main() {
	// Read CSV data
	citiesFile, err := os.Open("data/cities15000.txt")
	if err != nil {
		log.Fatalf("Opening file failed: %v", err)
	}
	r := csv.NewReader(citiesFile)
	r.Comma = '\t'
	r.Comment = '#'

	// Track collisions
	collisions := make(map[string]bool)

	// Pick out useful information
	cities := make(map[string]City)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to read CSV: %v", err)
		}
		key, city := CityFromRecord(record)

		// TODO: Reimplement collision rewriter

		// Remap collisions
		if _, ok := collisions[key]; ok {
			key = key + "_" + city.Admin1Code + "_" + city.CountryCode
		}

		// Check for collisions
		if existing, ok := cities[key]; ok {
			if existing.CountryCode == city.CountryCode {
				log.Printf("Warning: Repeat entry with same country code for %s (please compare %s with %s)", key, city.Timezone, existing.Timezone)
			} else if existing.Timezone == city.Timezone {
				log.Printf("Warning: Repeat entry with same timezone for %s", key)
			} else {
				log.Printf("Warning: Collision entry found for %s. Rewriting (%s but there is %s)", key, city.CountryCode, existing.CountryCode)
				cities[key+"_"+existing.Admin1Code+"_"+existing.CountryCode] = existing
				delete(cities, "key")
				collisions[key] = true
			}
		}
		cities[key] = city
	}

	// Group data
	data := Data{
		Cities: cities,
	}

	// Encode JSON file
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}

	// Write JSON file
	err = ioutil.WriteFile("js/data.json", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}
}

type Data struct {
	Cities map[string]City
}

type City struct {
	Names       []string `json:"n"`
	Admin1Code  string   `json:"a"`
	CountryCode string   `json:"c"`
	Timezone    string   `json:"t"`
}

// TODO: might be better to have IDs be City_Administrative_SG and the City struct having Names, Administrative, Country as pure text

func CityFromRecord(record []string) (string, City) {
	name := normalizeName(record[2])
	names := splitNames(record[3])
	admin1Code := record[10]
	countryCode := record[8]
	timezone := record[17]
	return name, City{
		Names:       names,
		Admin1Code:  admin1Code,
		CountryCode: countryCode,
		Timezone:    timezone,
	}
}

var re = regexp.MustCompile(`[^a-zA-Z1-9]`)

func normalizeName(name string) string {
	return re.ReplaceAllString(name, "_")
}

func splitNames(names string) []string {
	return strings.Split(names, ",")
}
*/
