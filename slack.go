package main

import (
	"fmt"
	// cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/nlopes/slack"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

// Fetch coin data from CMC and push it to slack.
func (s *SlackListener) pushCoinInfo(ticker string, channel string) error {

	// Grab the latest coin data from storage.
	coinInfo, err := getTicker(ticker)
	if err != nil {
		fmt.Errorf("can not obtain ticker from storage: %v", err)
	}

	attachment := slack.Attachment{
		Text: "Showing data for " + coinInfo.ID,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "24h Change",
				Value: fmt.Sprintf("%.2f%%", coinInfo.PercentChange24h),
			},
			slack.AttachmentField{
				Title: "Current value",
				Value: fmt.Sprintf("$%.2f", coinInfo.PriceUsd),
			},
		},
		Color: areWeHappy(coinInfo.PercentChange1h),
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(channel, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
