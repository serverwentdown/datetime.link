/*
This script reads in GeoNames data and creates a table of IDs to city names and
timezones. The IDs are created from the ASCII city name, with administrative
division level 1 name and country code as disambiguation. Examples of IDs are:

	Singapore-SG
	Ban_Bueng-Chon_Buri-TH
	Ashland-Oregon-US
	Ashland-California-US

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
	"sort"
	"strconv"
	"strings"

	"github.com/hbollon/go-edlib"
)

var regexName = regexp.MustCompile(`[^a-zA-Z1-9]+`)

// Data is all the data needed to map cities to timezones
type Data struct {
	Cities map[string]*City
}

// City represents a city that belongs inside an administrative division level 1
// and a country
type City struct {
	// Ref is the ASCII name of the city
	Ref string `json:"-"`
	// Name is the full UTF-8 name of the city
	Name           string   `json:"n"`
	AlternateNames []string `json:"an"`
	Timezone       string   `json:"t"`

	Population uint64 `json:"p"`

	Admin1  Admin1  `json:"a1"`
	Country Country `json:"c"`
}

// Admin1 represents an administrative division level 1
type Admin1 struct {
	// Code is the administrative division level 1 identifier, usually ISO-3166
	Code string `json:"-"`
	// Ref is the ASCII name of the administrative division level 1
	Ref string `json:"-"`
	// Name is the full UTF-8 name of the division
	Name string `json:"n"`
}

// Country represents a country
type Country struct {
	// CountryRef is the ISO-3166 2-letter country code
	Ref string `json:"-"`
	// Name is the full UTF-8 name of the country
	Name string `json:"n"`
}

func normalizeName(name string) string {
	simple := regexName.ReplaceAllString(name, "_")
	trimmed := strings.Trim(simple, "_")
	return trimmed
}

func splitNames(names string) []string {
	return strings.Split(names, ",")
}

type stringLengthSort []string

func (p stringLengthSort) Len() int           { return len(p) }
func (p stringLengthSort) Less(i, j int) bool { return len(p[i]) > len(p[j]) }
func (p stringLengthSort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func limitNames(primaryName string, names []string) ([]string, error) {
	sort.Sort(stringLengthSort(names))
	r := make([]string, 0, len(names))
	for _, n := range names {
		if n == primaryName || len(n) <= 0 {
			continue
		}
		// Skip abbreviation-like names
		if strings.ToUpper(n) == n {
			continue
		}
		// Skip almost the same names
		res, err := edlib.FuzzySearchThreshold(n, r, 0.82, edlib.Levenshtein)
		if err != nil {
			return nil, err
		}
		if len(res) != 0 {
			continue
		}
		// Skip substrings
		skipSubstr := false
		for _, longer := range r {
			if strings.HasPrefix(longer, n) {
				skipSubstr = true
			}
		}
		if skipSubstr {
			continue
		}
		// Limit
		if len(r) > 10 {
			break
		}
		r = append(r, n)
	}
	return r, nil
}

func extendRef(refs ...string) string {
	return strings.Join(refs, "-")
}

func readAdmin1Divisions(f string) (map[string]Admin1, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'

	m := make(map[string]Admin1)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		code := record[0]
		ref := normalizeName(record[2])
		name := record[1]
		m[code] = Admin1{
			Code: code,
			Ref:  ref,
			Name: name,
		}

	}
	return m, nil
}

func readCountries(f string) (map[string]Country, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'

	m := make(map[string]Country)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		ref := record[0]
		name := record[4]
		m[ref] = Country{
			Ref:  ref,
			Name: name,
		}
	}
	return m, nil
}

func readCities(f string, countries map[string]Country, admin1s map[string]Admin1) (map[string]*City, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'

	m := make(map[string]*City)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		name := record[1]
		ref := normalizeName(record[2])
		alternateNames, err := limitNames(name, splitNames(record[3]))
		if err != nil {
			return nil, err
		}
		admin1Code := record[10]
		countryRef := record[8]
		population, err := strconv.ParseUint(record[14], 10, 64)
		if err != nil {
			return nil, err
		}
		timezone := record[17]

		// Resolve Country and Admin1
		country := countries[countryRef]
		admin1 := admin1s[countryRef+"."+admin1Code]

		// Bulid a full formed ID
		eref := extendRef(ref, admin1.Ref, country.Ref)
		if len(admin1.Ref) <= 0 {
			eref = extendRef(ref, country.Ref)
		}

		c := &City{
			Ref:            ref,
			Name:           name,
			AlternateNames: alternateNames,
			Timezone:       timezone,
			Population:     population,
			Admin1:         admin1,
			Country:        country,
		}

		// Warn if there exists a similar city
		if e, ok := m[eref]; ok {
			if !(e.Ref == c.Ref && e.Admin1.Ref == c.Admin1.Ref && e.Country.Ref == e.Country.Ref) {

				log.Printf("WARNING: existing city %s: %v %v", eref, c, e)
			}
		}

		m[eref] = c
	}
	return m, nil
}

func main() {
	admin1s, err := readAdmin1Divisions("../data/admin1CodesASCII.txt")
	if err != nil {
		log.Fatalf("Reading administrative divisions level 1 failed")
	}
	countries, err := readCountries("../data/countryInfo.txt")
	if err != nil {
		log.Fatalf("Reading countries failed")
	}
	cities, err := readCities("../data/cities15000.txt", countries, admin1s)
	if err != nil {
		log.Fatalf("Reading cities failed")
	}

	// Group data
	data := Data{
		Cities: cities,
	}
	// Encode JSON file
	//b, err := json.MarshalIndent(data, " ", " ")
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
	// Write JSON file
	err = ioutil.WriteFile("../js/data.json", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}
}
