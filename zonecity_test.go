package main

import (
	"testing"

	"github.com/serverwentdown/datetime.link/data"
)

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
