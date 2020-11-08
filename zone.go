package main

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/serverwentdown/datetime.link/data"
)

// ErrZoneNotFound is thrown when a zone string has no match
var ErrZoneNotFound = errors.New("zone not found")

// ErrZoneOffsetInvalid is thrown when a zone string is an invalid zone offset
var ErrZoneOffsetInvalid = errors.New("offset zone invalid")

var zoneOffsetRegexp = regexp.MustCompile(`^[+-][0-9]{2}:[0-9]{2}$`)

// SearchCities looks up a city by it's reference
func SearchCities(cities map[string]*data.City, city string) (*data.City, error) {
	// For now, simple map read will do
	if city, ok := cities[city]; ok {
		return city, nil
	}
	return nil, ErrZoneNotFound
}

// ParseZoneOffset parses a zone string into a time.Location
func ParseZoneOffset(zone string) (*time.Location, error) {
	if !zoneOffsetRegexp.MatchString(zone) {
		return nil, ErrZoneOffsetInvalid
	}
	// Assume that if it satisfies the regex, it satisfies the length and won't
	// fail to parse
	d := 0
	if zone[0] == '+' {
		d = 1
	}
	if zone[0] == '-' {
		d = -1
	}
	h, _ := strconv.ParseUint(zone[1:1+2], 10, 64)
	// Allow hour offsets greater that 24
	m, _ := strconv.ParseUint(zone[1+3:1+3+2], 10, 64)
	if m >= 60 {
		return nil, ErrZoneOffsetInvalid
	}
	return time.FixedZone("UTC"+zone, d*(int(h)*60*60+int(m)*60)), nil
}

// Zone represents any form of zone offset: this could be a city or a fixed
// offset from UTC
type Zone struct {
	// City is optional
	City *data.City
	// Offset is optional
	Offset *time.Location
}

// ResolveZone resolves a zone string into a Zone
func ResolveZone(cities map[string]*data.City, zone string) (Zone, error) {
	// Try parsing as a zone offset
	offset, err := ParseZoneOffset(zone)
	if err == nil {
		return Zone{
			Offset: offset,
		}, nil
	}
	// Parse as a city
	city, err := SearchCities(cities, zone)
	if err == nil {
		return Zone{
			City: city,
		}, nil
	}
	return Zone{}, err
}
