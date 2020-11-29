package main

import (
	"sort"
	"strings"

	"github.com/hbollon/go-edlib"
	"github.com/serverwentdown/datetime.link/data"
)

// FullSearchCities uses a very basic iterative method to search for cities
// with the given string
func FullSearchCities(cities map[string]*data.City, zone string) ([]*data.City, error) {
	// TODO: optimisations
	ratings := []cityRatings{}
	for _, city := range cities {
		rating, err := compareCity(city, zone)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, cityRatings{city, rating})
	}
	sort.Slice(ratings, func(i, j int) bool { return ratings[i].Rating > ratings[j].Rating })
	topCities := make([]*data.City, 10)
	for i := 0; i < 10; i++ {
		topCities[i] = ratings[i].City
		//l.Debug("city", zap.String("n", topCities[i].Name), zap.Float32("r", ratings[i].Rating))
	}
	return topCities, nil
}

type cityRatings struct {
	City   *data.City
	Rating float32
}

func compareCity(city *data.City, zone string) (float32, error) {
	// City Name is preferred
	cityDistance, err := compare(city.Name, zone)
	if err != nil {
		return 0, err
	}
	for _, altname := range city.AlternateNames {
		altnameDistance, err := compare(altname, zone)
		if err != nil {
			return 0, err
		}
		altnameDistance *= 0.9
		cityDistance = floatMax(cityDistance, altnameDistance)
	}
	// Admin1 Name is next preferred
	admin1Distance, err := compare(city.Admin1.Name, zone)
	if err != nil {
		return 0, err
	}
	// Country Name is next preferred
	countryDistance, err := compare(city.Country.Name, zone)
	if err != nil {
		return 0, err
	}
	// Merge 3 values
	rating := floatMax(cityDistance, admin1Distance*0.9, countryDistance*0.9)
	return rating, nil
}

func compare(str1, str2 string) (float32, error) {
	algo := edlib.JaroWinkler
	//algo := edlib.Levenshtein
	res, err := edlib.StringsSimilarity(strings.ToLower(str1), strings.ToLower(str2), algo)
	return res * res * res, err
}

func floatMax(a float32, bs ...float32) float32 {
	for _, b := range bs {
		if a > b {
			continue
		}
		a = b
	}
	return a
}
