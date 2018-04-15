package main

import (
	"time"
)

// The Ticker struct holds a map of Coin structs, the last time
// data was fetched successfully and in the case of an error, the error message
type Ticker struct {
	Coins      map[string]*Coin
	LastUpdate time.Time
	Error      error
}

// The Coin struct holds the data pertaining to a specific coin
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int16   `json:"rank,string"`
	PriceUsd         float64 `json:"price_usd,string"`
	PriceBtc         float64 `json:"price_btc,string"`
	Usd24hVolume     float64 `json:"24h_volume_usd,string"`
	MarketCapUsd     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1h  float64 `json:"percent_change_1h,string"`
	PercentChange24h float64 `json:"percent_change_24h,string"`
	PercentChange7d  float64 `json:"percent_change_7d,string"`
	LastUpdated      string  `json:"last_updated"`
}
