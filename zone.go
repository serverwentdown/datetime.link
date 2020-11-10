package main

import (
	"time"

	"github.com/serverwentdown/datetime.link/data"
	"go.uber.org/zap"
)

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

// FirstName returns the primary name of the zone
func (z Zone) FirstName() string {
	if z.IsOffset() {
		return z.Offset.String()
	} else if z.City != nil {
		if len(z.City.Admin1.Name) > 0 {
			return z.City.Name + ", " + z.City.Admin1.Name
		}
		return z.City.Name
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
