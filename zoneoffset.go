package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrZoneOffsetInvalid is thrown when a zone string is an invalid zone offset
var ErrZoneOffsetInvalid = errors.New("offset zone invalid")

const zoneHour = 60 * 60
const zoneMinute = 60

var zoneOffsetRegexp = regexp.MustCompile("^[+-\u2212][0-9]{2}:[0-9]{2}$")

// ParseZoneOffset parses a zone string into an offset
func ParseZoneOffset(zone string) (int, error) {
	if !zoneOffsetRegexp.MatchString(zone) {
		return 0, ErrZoneOffsetInvalid
	}
	// Assume that if it satisfies the regex, it satisfies the length and won't
	// fail to parse
	d := 0
	if strings.HasPrefix(zone, "+") {
		d = 1
	}
	if strings.HasPrefix(zone, "-") || strings.HasPrefix(zone, "\u2212") {
		d = -1
	}
	zone = strings.TrimLeft(zone, "+-\u2212")
	parts := strings.Split(zone, ":")
	h, _ := strconv.ParseUint(parts[0], 10, 64)
	// Allow hour offsets greater that 24
	m, _ := strconv.ParseUint(parts[1], 10, 64)
	if m >= 60 {
		return 0, ErrZoneOffsetInvalid
	}
	offset := d * (int(h)*zoneHour + int(m)*zoneMinute)
	return offset, nil
}

// FormatZoneOffset formats an offset into a string
func FormatZoneOffset(offset int) string {
	neg := offset < 0
	s := '+'
	if neg {
		s = '\u2212'
		offset = -offset
	}
	if offset == 0 {
		return "\u00b100:00"
	}
	h := offset / zoneHour
	m := (offset % zoneHour) / zoneMinute
	return fmt.Sprintf("%c%02d:%02d", s, h, m)
}
