package main

import (
	"fmt"
	"log"
	"strings"

	cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/nlopes/slack"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {

		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			info := rtm.GetInfo()
			if ev.User != info.User.ID {
				if err := s.handleMessageEvent(ev); err != nil {
					log.Printf("[ERROR] Failed to handle message: %s", err)
				}
			}
		}
	}

}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Do not react to itself.
	if ev.User == s.botID {
		return nil
	}
	if !strings.HasPrefix(ev.Msg.Text, "<@"+s.botID+">") {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]

	// Allowed commands (coins)
	commands := []string{"eth", "ether", "ethereum", "btc", "bitcoin", "ltc", "litecoin"}

	ack := false
	for _, cmd := range commands {
		if m[0] == cmd {
			ack = true
			break
		}
	}

	if !ack {
		return fmt.Errorf("invalid message")
	}

	var tracker string
	if m[0] == "eth" || m[0] == "ether" || m[0] == "ethereum" {
		tracker = "ethereum"
	}
	if m[0] == "btc" || m[0] == "bitcoin" {
		tracker = "bitcoin"
	}
	if m[0] == "ltc" || m[0] == "litecoin" {
		tracker = "litecoin"
	}

	// Fetch coin data from CMC
	coinInfo, err := cmc.GetCoinData(tracker)
	if err != nil {
		fmt.Errorf("coin data retrieval failed: %v", err)
	}

	attachment := slack.Attachment{
		Text: "Showing data for " + coinInfo.ID,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "24h Change",
				Value: fmt.Sprintf("%.2f%%", coinInfo.PercentChange24H),
			},
			slack.AttachmentField{
				Title: "Current value",
				Value: fmt.Sprintf("$%.2f", coinInfo.PriceUSD),
			},
		},
		Color: areWeHappy(coinInfo.PercentChange1H),
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}

func areWeHappy(change float64) string {
	if change < 0 {
		return "warning"
	}
	return "good"
}
