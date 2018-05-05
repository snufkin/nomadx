package main

import (
	"testing"
)

var testHappyString = []struct {
	in      float64
	out     string
	success bool
}{
	{-0.1, slackWarning, true},
	{-0.1, slackGood, false},
	{0, slackWarning, false},
	{0, slackGood, true},
	{0.5, slackWarning, false},
	{0.5, slackGood, true},
}

func TestAreWeHappy(t *testing.T) {
	for _, testChange := range testHappyString {
		output := areWeHappy(testChange.in)
		if success := (output == testChange.out); success != testChange.success {
			t.Error("For", testChange.in, "expected", testChange.out, "got", output)
		}
	}
}
