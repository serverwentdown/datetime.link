package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/serverwentdown/datetime.link/data"
	"go.uber.org/zap"
)

// ErrZoneNotFound is thrown when a zone string has no match
var ErrZoneNotFound = errors.New("zone not found")

// ErrZoneOffsetInvalid is thrown when a zone string is an invalid zone offset
var ErrZoneOffsetInvalid = errors.New("offset zone invalid")

const zoneHour = 60 * 60
const zoneMinute = 60

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
func ParseZoneOffset(zone string) (int, error) {
	if !zoneOffsetRegexp.MatchString(zone) {
		return 0, ErrZoneOffsetInvalid
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
		return 0, ErrZoneOffsetInvalid
	}
	offset := d * (int(h)*zoneHour + int(m)*zoneMinute)
	return offset, nil
}

// FormatZoneOffset formats an offset into a string
func FormatZoneOffset(offset int) string {
	neg := offset < 0
	s := '+'
	if neg {
		s = '-'
		offset = -offset
	}
	if offset == 0 {
		return "\u00B100:00"
	}
	h := offset / zoneHour
	m := (offset % zoneHour) / zoneMinute
	return fmt.Sprintf("%c%02d:%02d", s, h, m)
}

// Zone represents any form of zone offset: this could be a city or a fixed
// offset from UTC
type Zone struct {
	// City is optional
	City *data.City
	// Offset is optional
	Offset *time.Location
}

// IsOffset is true when the Zone is an offset instead of a city
func (z Zone) IsOffset() bool {
	return z.Offset != nil
}

// Name returns the name of the zone
func (z Zone) Name() string {
	if z.IsOffset() {
		return z.Offset.String()
	} else if z.City != nil {
		return z.City.FullName()
	}
	return ""
}

// Location returns the time.Location of the zone. Useful for other functions
func (z Zone) Location() *time.Location {
	if z.IsOffset() {
		return z.Offset
	} else if z.City != nil {
		loc, err := time.LoadLocation(z.City.Timezone)
		if err != nil {
			// This is a really bad situation
			// TODO: validate this at boot time
			l.Error("unable to find timezone", zap.String("timezone", z.City.Timezone))
		}
		return loc
	}
	return nil
}

// TimeOffset returns the timezone offset at a specific time
func (z Zone) TimeOffset(t time.Time) int {
	if l := z.Location(); l != nil {
		_, offset := t.In(l).Zone()
		return offset
	}
	return -1 // TODO: better invalid handling
}

// ResolveZone resolves a zone string into a Zone
func ResolveZone(cities map[string]*data.City, zone string) (Zone, error) {
	// Try parsing as a zone offset
	offset, err := ParseZoneOffset(zone)
	if err == nil {
		offsetZone := time.FixedZone("UTC "+FormatZoneOffset(offset), offset)
		return Zone{
			Offset: offsetZone,
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
