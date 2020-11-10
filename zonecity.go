package main

import (
	"errors"

	"github.com/serverwentdown/datetime.link/data"
)

// ErrZoneNotFound is thrown when a zone string has no match
var ErrZoneNotFound = errors.New("zone not found")

// SearchCities looks up a city by it's reference
func SearchCities(cities map[string]*data.City, city string) (*data.City, error) {
	// For now, simple map read will do
	if city, ok := cities[city]; ok {
		return city, nil
	}
	return nil, ErrZoneNotFound
}
