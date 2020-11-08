package data

import (
	"encoding/json"
	"io/ioutil"
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
