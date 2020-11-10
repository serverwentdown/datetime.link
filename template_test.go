package main

import (
	"testing"
)

func TestTemplateFuncFormatOffset(t *testing.T) {
	want, got := "+06:06", templateFuncFormatOffset(6*60*60+6*60)
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "-12:15", templateFuncFormatOffset(-(12*60*60 + 15*60))
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "\u00B100:00", templateFuncFormatOffset(-(0*60*60 + 0*60))
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}

	want, got = "+00:01", templateFuncFormatOffset(0*60*60+1*60)
	if want != got {
		t.Fatalf("got offset %v, want offset %v", got, want)
	}
}
