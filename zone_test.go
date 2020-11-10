package main

import (
	"testing"
	"time"

	"github.com/serverwentdown/datetime.link/data"
)

func TestParseZoneOffset(t *testing.T) {
	loc, err := ParseZoneOffset("+08:00")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		time := time.Date(2020, time.November, 8, 23, 9, 0, 0, loc).Unix()
		want := int64(1604848140)
		if time != want {
			t.Errorf("got %d, want %d", time, want)
		}
	}

	loc, err = ParseZoneOffset("-00:30")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		time := time.Date(2020, time.November, 8, 14, 39, 0, 0, loc).Unix()
		want := int64(1604848140)
		if time != want {
			t.Errorf("got %d, want %d", time, want)
		}
	}

	loc, err = ParseZoneOffset("-0030")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	loc, err = ParseZoneOffset("00:30")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	loc, err = ParseZoneOffset("+08:60")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	loc, err = ParseZoneOffset("+08:-6")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}
}

func TestSearchCities(t *testing.T) {
	cities, err := data.ReadCities()
	if err != nil {
		panic(err)
	}

	city, err := SearchCities(cities, "Singapore-SG")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	}
	wantName := "Singapore"
	wantZone := "Asia/Singapore"
	if city.Name != wantName || city.Timezone != wantZone {
		t.Errorf("want %v %v, got %v", wantName, wantZone, city)
	}

	city, err = SearchCities(cities, "Yuzhno_Sakhalinsk-Sakhalin_Oblast-RU")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	}
	wantName = "Yuzhno-Sakhalinsk"
	wantZone = "Asia/Sakhalin"
	if city.Name != wantName || city.Timezone != wantZone {
		t.Errorf("want %v %v, got %v", wantName, wantZone, city)
	}

	_, err = SearchCities(cities, "Nowhere")
	if err != ErrZoneNotFound {
		t.Errorf("want error %v, got error %v", ErrZoneNotFound, err)
	}
}

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
	wantOffset, _ := ParseZoneOffset("+04:00")
	if zone.Offset.String() != wantOffset.String() {
		t.Errorf("want Offset %v, got Offset %v", wantOffset, zone.Offset)
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

func BenchmarkReadCities(b *testing.B) {
	// This does take quite a while
	for i := 0; i < b.N; i++ {
		_, _ = data.ReadCities()
	}
}

func BenchmarkSearchCities(b *testing.B) {
	cities, err := data.ReadCities()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SearchCities(cities, "Yuzhno_Sakhalinsk-Sakhalin_Oblast-RU")
	}
}
