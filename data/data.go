package data

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// ReadCities opens the file "data/cities.json" and reads it into a map
func ReadCities() (map[string]*City, error) {
	cities := make(map[string]*City)

	buf, err := ioutil.ReadFile("data/cities.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func extendName(names ...string) string {
	return strings.Join(names, ", ")
}

// FullName returns a fully qualified human readable name
func (c City) FullName() string {
	if len(c.Admin1.Name) > 0 {
		return extendName(c.Name, c.Admin1.Name, c.Country.Ref)
	}
	return extendName(c.Name, c.Country.Ref)
}
