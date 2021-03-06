package main

import (
	"fmt"
)

const slackWarning string = "warning"
const slackGood string = "good"

// Determine the slack colour string based on the sign of the number.
func areWeHappy(change float64) string {
	if change < 0 {
		return slackWarning
	}
	return slackGood
}

// Generate a Slack formatted string for a given timestamp.
// <!date^unix_epoch_timestamp^string_containing_date_tokens^optional_link|fallback_text>
func slackTime(timeStamp int64, slackFormat string) string {
	return fmt.Sprintf("<!date^%d^%v|generated at %d>", timeStamp, slackFormat, timeStamp)
}
