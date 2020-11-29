package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ErrComponentsMismatch is thrown when the URL has empty or missing components
var ErrComponentsMismatch = errors.New("missing or too many path components")

// ErrInvalidTime is thrown in a time.ParseError
var ErrInvalidTime = errors.New("invalid ISO 8601 time")

var timeRFC3339NoSec = "2006-01-02T15:04Z07:00"
var timeFormats = []string{time.RFC3339, timeRFC3339NoSec}

// Request is a parsed datetime URL
type Request struct {
	Time  time.Time
	Zones []string
}

// ParseRequest parses an input URL into a Request
func ParseRequest(u *url.URL) (Request, error) {
	var err error

	parts := strings.Split(u.Path, "/")[1:]
	if len(parts) > 2 || len(parts) < 1 {
		return Request{}, ErrComponentsMismatch
	}

	// Parse time portion
	var t time.Time
	timeString := parts[0]
	if len(timeString) == 0 {
		return Request{}, ErrComponentsMismatch
	}
	if timeString == "now" {
		t = time.Now()
	} else {
		for _, f := range timeFormats {
			t, err = time.Parse(f, timeString)
			if err == nil {
				break
			}
		}
		if err != nil {
			return Request{}, fmt.Errorf("%w: %v", ErrInvalidTime, err)
		}
	}

	// Split zones
	var z []string
	zoneString := "local"
	if len(parts) >= 2 {
		zoneString = parts[1]
	}
	if len(zoneString) == 0 {
		return Request{}, ErrComponentsMismatch
	}
	z = strings.Split(zoneString, ",")

	return Request{
		Time:  t,
		Zones: z,
	}, nil
}
