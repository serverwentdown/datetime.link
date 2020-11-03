package main

import (
	"testing"
)

func TestChooseResponseType(t *testing.T) {
	r := chooseResponseType("text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	if r != responseHTML {
		t.Errorf("expecting html, got %v", r)
	}
	r = chooseResponseType("application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	if r != responseAny {
		t.Errorf("expecting html, got %v", r)
	}
}
