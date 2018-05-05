package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strings"
)

// Initialise the global storage of coins.
var Storage = TickerStorage{}

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

// Generic message parser and delegator function.
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

	syn, err := getTickerSynonym(m[0])
	if err {
		return err
	}

	if err := s.pushCoinInfo(syn, ev.Channel); err != nil {
		return fmt.Errorf("coin info push failed: %s", err)
	}

	return nil
}
