package main

import (
	"strconv"
	"strings"
)

// Snippet from https://github.com/emicklei/go-restful/blob/master/mime.go
// MIT License

type mime struct {
	media   string
	quality float64
}

// insertMime adds a mime to a list and keeps it sorted by quality.
func insertMime(l []mime, e mime) []mime {
	for i, each := range l {
		// if current mime has lower quality then insert before
		if e.quality > each.quality {
			left := append([]mime{}, l[0:i]...)
			return append(append(left, e), l[i:]...)
		}
	}
	return append(l, e)
}

const qFactorWeightingKey = "q"

// sortedMimes returns a list of mime sorted (desc) by its specified quality.
// e.g. text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3
func sortedMimes(accept string) (sorted []mime) {
	for _, each := range strings.Split(accept, ",") {
		typeAndQuality := strings.Split(strings.Trim(each, " "), ";")
		if len(typeAndQuality) == 1 {
			sorted = insertMime(sorted, mime{typeAndQuality[0], 1.0})
		} else {
			// take factor
			qAndWeight := strings.Split(typeAndQuality[1], "=")
			if len(qAndWeight) == 2 && strings.Trim(qAndWeight[0], " ") == qFactorWeightingKey {
				f, err := strconv.ParseFloat(qAndWeight[1], 64)
				if err != nil {
					// do nothing
				} else {
					sorted = insertMime(sorted, mime{typeAndQuality[0], f})
				}
			} else {
				sorted = insertMime(sorted, mime{typeAndQuality[0], 1.0})
			}
		}
	}
	return
}

type responseType int

const (
	responsePlain   responseType = iota
	responseHTML    responseType = iota
	responseJSON    responseType = iota
	responseAny     responseType = iota
	responseUnknown responseType = iota
)

const (
	responsePlainMime = "text/plain"
	responseHTMLMime  = "text/html"
	responseJSONMime  = "application/json"
	responseAnyMime   = "*/*"
)

// chooseResponse returns a response type from an accept header
func chooseResponseType(accept string) responseType {
	acceptSorted := sortedMimes(accept)
	for _, m := range acceptSorted {
		if m.media == responsePlainMime {
			return responsePlain
		}
		if m.media == responseHTMLMime {
			return responseHTML
		}
		if m.media == responseJSONMime {
			return responseJSON
		}
		if m.media == responseAnyMime {
			return responseAny
		}
	}
	return responseUnknown
}
