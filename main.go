package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

func main() {
	var c BotConfig

	if err := envconfig.Process("nomadx", &c); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("[INFO] Start slack event listening")
	client := slack.New(c.BotToken)

	slackListener := &SlackListener{
		client:    client,
		botID:     c.BotID,
		channelID: c.ChannelID,
	}

	period := 300
	tickers := []string{"ethereum", "litecoin", "bitcoin"}
	go slackListener.fetchData(period, tickers)
	slackListener.ListenAndResponse()
}
