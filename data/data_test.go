package data

import (
	"testing"
)

func TestCityFullName(t *testing.T) {
	c := City{
		Name: "'cit-y N'ame-",
		Admin1: Admin1{
			Name: "a-dmin1, Name'",
		},
		Country: Country{
			Name: "Country Name",
			Ref:  "CN",
		},
	}

	got, want := c.FullName(), "'cit-y N'ame-, a-dmin1, Name', CN"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	c = City{
		Name:   "City Name",
		Admin1: Admin1{},
		Country: Country{
			Name: "Name Country",
			Ref:  "CN",
		},
	}

	got, want = c.FullName(), "City Name, CN"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
