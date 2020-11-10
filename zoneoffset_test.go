package main

import (
	"testing"
	"time"
)

func TestParseZoneOffset(t *testing.T) {
	offset, err := ParseZoneOffset("+08:00")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		want := 8*60*60 + 0*60
		if offset != want {
			t.Errorf("got %d, want %d", offset, want)
		}
	}

	offset, err = ParseZoneOffset("-01:30")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		want := -(1*60*60 + 30*60)
		if offset != want {
			t.Errorf("got %d, want %d", offset, want)
		}
	}

	offset, err = ParseZoneOffset("\u221251:59")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		want := -(51*60*60 + 59*60)
		if offset != want {
			t.Errorf("got %d, want %d", offset, want)
		}
	}

	_, err = ParseZoneOffset("-0030")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	_, err = ParseZoneOffset("00:30")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	_, err = ParseZoneOffset("+08:60")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	_, err = ParseZoneOffset("+08:-6")
	if err != ErrZoneOffsetInvalid {
		t.Errorf("want error %v, got error %v", ErrZoneOffsetInvalid, err)
	}

	offset, err = ParseZoneOffset("+08:00")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		loc := time.FixedZone("UTC "+FormatZoneOffset(offset), offset)
		time := time.Date(2020, time.November, 8, 23, 9, 0, 0, loc).Unix()
		want := int64(1604848140)
		if time != want {
			t.Errorf("got %d, want %d", time, want)
		}
	}

	offset, err = ParseZoneOffset("-00:30")
	if err != nil {
		t.Errorf("want error %v, got error %v", nil, err)
	} else {
		loc := time.FixedZone("UTC "+FormatZoneOffset(offset), offset)
		time := time.Date(2020, time.November, 8, 14, 39, 0, 0, loc).Unix()
		want := int64(1604848140)
		if time != want {
			t.Errorf("got %d, want %d", time, want)
		}
	}
}

func TestFormatZoneOffset(t *testing.T) {
	want, got := "+06:06", FormatZoneOffset(6*60*60+6*60)
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "\u221212:15", FormatZoneOffset(-(12*60*60 + 15*60))
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "\u00B100:00", FormatZoneOffset(-(0*60*60 + 0*60))
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "+00:01", FormatZoneOffset(0*60*60+1*60)
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}
}
