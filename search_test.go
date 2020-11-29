package main

import (
	"log"
	"testing"

	"github.com/hbollon/go-edlib"
	"github.com/serverwentdown/datetime.link/data"
)

func TestEditDistance(t *testing.T) {
	res, err := edlib.StringsSimilarity("Singapore", "Sing", edlib.JaroWinkler)
	if err != nil {
		return
	}
	log.Printf("%f", res)
}

func BenchmarkCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = compare("Random String That Is Quite Long", "Singapore")
	}
}

func BenchmarkCompareCity(b *testing.B) {
	cities, err := data.ReadCities()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compareCity(cities["Singapore-SG"], "Singapore")
	}
}

func BenchmarkFullSearchCities(b *testing.B) {
	cities, err := data.ReadCities()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FullSearchCities(cities, "Singapore")
	}
}
