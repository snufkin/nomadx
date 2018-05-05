package main

import (
	"testing"
)

var testTickerSynonyms = []struct {
	in      string
	out     string
	success bool
}{
	{"ethereum", "eth", true},
	{"eth", "eth", true},
	{"litecoin", "ltc", true},
	{"ltc", "ltc", true},
	{"bitcoin", "btc", true},
	{"btc", "btc", true},
	{"bcash", "btc", false},
}

func TestSynonyms(t *testing.T) {
	for _, testC := range testTickerSynonyms {
		syn, _ := getTickerSynonym(testC.in) // TODO check the error
		if success := (syn == testC.out); success != testC.success {
			t.Error("For", testC.in, "expected", testC.out, "got", syn)
		}
	}
}
