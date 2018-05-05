package main

import (
	"fmt"
)

var currencyNames = map[string]string{
	"ethereum": "eth",
	"eth":      "eth",
	"litecoin": "ltc",
	"ltc":      "ltc",
	"bitcoin":  "btc",
	"btc":      "btc",
}

func getTickerSynonym(synonym string) (ticker string, err error) {
	if ticker, ok := currencyNames[synonym]; ok {
		return ticker, nil
	} else {
		return "", fmt.Errorf("synonyms: unknown ticker requested: %s", synonym)
	}
}
