package main

import (
	cmc "github.com/miguelmota/go-coinmarketcap"
	"time"
)

type BotConfig struct {
	BotToken  string `required:"true"`
	ChannelID string `required:"true"`
	BotID     string `required:"true"`
}

// The Ticker struct holds a map of Coin structs, the last time
// data was fetched successfully and in the case of an error, the error message
type TickerStorage struct {
	Coins       map[string][]cmc.Coin
	LastUpdate  time.Time
	Error       error
	LastUpdates map[string]time.Time
}
