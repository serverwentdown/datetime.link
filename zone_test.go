package main

import (
	"testing"
	"time"

	"github.com/serverwentdown/datetime.link/data"
)

func TestResolveZone(t *testing.T) {
	cities, err := data.ReadCities()
	if err != nil {
		panic(err)
	}

	zone, err := ResolveZone(cities, "Singapore-SG")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	}
	wantCity, _ := SearchCities(cities, "Singapore-SG")
	// TODO: instead of pointer comparison, do .Equals()
	if zone.City != wantCity {
		t.Errorf("want City %v, got City %v", wantCity, zone.City)
	}
	if zone.Offset != nil {
		t.Errorf("want Offset %v, got Offset %v", nil, zone.Offset)
	}

	zone, err = ResolveZone(cities, "+04:00")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	}
	if zone.City != nil {
		t.Errorf("want City %v, got City %v", nil, zone.City)
	}
	gotTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, zone.Offset)
	wantTime := time.Date(2019, time.December, 31, 20, 0, 0, 0, time.UTC)
	if !gotTime.Equal(wantTime) {
		t.Errorf("want time %v, got time %v", wantTime, gotTime)
	}

	zone, err = ResolveZone(cities, "+04:80")
	if err != ErrZoneNotFound {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	zone, err = ResolveZone(cities, "04:80")
	if err != ErrZoneNotFound {
		t.Errorf("want error %v, got error %v", ErrZoneNotFound, err)
	}
}
