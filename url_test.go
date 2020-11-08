package main

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mustURLParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func TestURLParse(t *testing.T) {
	u := mustURLParse("http://test/2020-06-02T14:00+08:00/Singapore,Malaysia")
	got, err := ParseRequest(u)
	if err != nil {
		t.Errorf("mismatch: got error %v", err)
		return
	}
	want := Request{
		time.Date(2020, 6, 2, 14, 0, 0, 0, time.FixedZone("UTC +8", 8*60*60)),
		[]string{"Singapore", "Malaysia"},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("mismatch: \n%v", cmp.Diff(got, want))
	}

	u = mustURLParse("http://test/2019-04-30T18:00:00Z/Nowhere")
	got, err = ParseRequest(u)
	if err != nil {
		t.Errorf("mismatch: got error %v", err)
		return
	}
	want = Request{
		time.Date(2019, 4, 30, 18, 0, 0, 0, time.FixedZone("UTC", 0)),
		[]string{"Nowhere"},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("mismatch: \n%v", cmp.Diff(got, want))
	}
}

func TestURLParseFail(t *testing.T) {
	u := mustURLParse("http://test/2002-08-30T14:00+06:00/")
	_, err := ParseRequest(u)
	if !errors.Is(err, ErrMissingComponent) {
		t.Errorf("mismatch: got error %v, want error %v", err, ErrMissingComponent)
		return
	}

	u = mustURLParse("http://test/")
	_, err = ParseRequest(u)
	if !errors.Is(err, ErrMissingComponent) {
		t.Errorf("mismatch: got error %v, want error %v", err, ErrMissingComponent)
		return
	}

	u = mustURLParse("http://test/2000-01-13T00:00Z08:00/hi")
	_, err = ParseRequest(u)
	_, isParseError := err.(*time.ParseError)
	if !isParseError {
		t.Errorf("mismatch: got error %v, want time.ParseError", err)
		return
	}

	u = mustURLParse("http://test/2000-01-13 00:00+08:00/hi")
	_, err = ParseRequest(u)
	_, isParseError = err.(*time.ParseError)
	if !isParseError {
		t.Errorf("mismatch: got error %v, want time.ParseError", err)
		return
	}

}
