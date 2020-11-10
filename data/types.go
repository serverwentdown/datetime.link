package data

// City represents a city that belongs inside an administrative division level 1
// and a country
type City struct {
	// Ref is the ASCII name of the city
	Ref string `json:"-"`
	// Name is the full UTF-8 name of the city
	Name           string   `json:"n"`
	AlternateNames []string `json:"an"`
	Timezone       string   `json:"t"`

	Population uint64 `json:"p"`

	Admin1  Admin1  `json:"a1"`
	Country Country `json:"c"`
}

// Admin1 represents an administrative division level 1
type Admin1 struct {
	// Code is the administrative division level 1 identifier, usually ISO-3166
	Code string `json:"-"`
	// Ref is the ASCII name of the administrative division level 1
	Ref string `json:"-"`
	// Name is the full UTF-8 name of the division
	Name string `json:"n"`
}

// Country represents a country
type Country struct {
	// CountryRef is the ISO-3166 2-letter country code
	Ref string `json:"r"`
	// Name is the full UTF-8 name of the country
	Name string `json:"n"`
}
