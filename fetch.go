package main

import (
	"fmt"
	cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/nlopes/slack"
	"time"
)

// Fetch data for a particular ticker from CMC
// Avoid invoking this directly as this pings the endpoint directly.
func fetchCMCTicker(ticker string) (cmc.Coin, error) {
	// Fetch coin data from CMC
	coinInfo, err := cmc.GetCoinData(ticker)
	if err != nil {
		fmt.Errorf("coin data retrieval failed: %v", err)
	}
	return coinInfo, err
}

// Ticker getter with caching wrapper (not less than 10 per minute)
func getTicker(ticker string) (coinInfo cmc.Coin, err error) {
	t := time.Now()

	// Do we even have the coin on record?
	coinList, ok := []cmc.Coin{}, false

	if coinList, ok = Storage.Coins[ticker]; ok {

		// We have the coin, and it is not stale.
		if !isTickerStale(Storage.LastUpdates[ticker], 6) {
			coinInfo = coinList[len(coinList)-1] // Last element of the entries.
		} else {
			// Stale coin, update storage.
			if coinInfo, err = fetchCMCTicker(ticker); err != nil {
				fmt.Errorf("ticker failed to refresh: %v", err)
			} else {
				Storage.Coins[ticker] = append(Storage.Coins[ticker], coinInfo)
				Storage.LastUpdates[ticker] = t
			}
		}
	} else {
		// No coin, need to fetch.
		if coinInfo, err = fetchCMCTicker(ticker); err != nil {
			Storage.Coins[ticker] = append(Storage.Coins[ticker], coinInfo)
		}
	}

	return
}

// Check if the ticker is stale based on the passed timeout. Missing means stale.
func isTickerStale(lastUpdated time.Time, timeLimit int64) bool {
	return time.Since(lastUpdated) > time.Duration(timeLimit)*time.Second
}

// Main loop to periodically download the status for all tickers.
func (s *SlackListener) fetchData(period int, coinTickers []string) {
	// Fetch coin data from CMC for all tickers
	ticker := time.NewTicker(time.Duration(period) * time.Second)

	for _ = range ticker.C {
		s.updateData(coinTickers)
	}
}

func (s *SlackListener) updateData(tickers []string) {
	t := time.Now()
	coinList := make([]cmc.Coin, 3)
	fmt.Printf("[%v] ", t.Format(time.Stamp))
	for i, tck := range tickers {
		coinInfo, err := cmc.GetCoinData(tck)
		if err != nil {
			fmt.Errorf("coin data retrieval failed: %v", err)
		}
		s.recordData(coinInfo)
		coinList[i] = coinInfo
	}
	s.postUpdate(coinList)
	fmt.Printf("\n")
}

func (s *SlackListener) recordData(coinInfo cmc.Coin) {
	fmt.Printf("%v@%.2f ", coinInfo.Symbol, coinInfo.PriceUSD)
}

func (s *SlackListener) postUpdate(coinList []cmc.Coin) error {
	t := time.Now()
	attachment := slack.Attachment{
		Text: fmt.Sprintf("[%v] hourly update", slackTime(t.Unix(), "{time}")),

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: coinList[0].ID,
				Value: fmt.Sprintf("$%.2f", coinList[0].PriceUSD),
				Short: true,
			},
			slack.AttachmentField{
				Title: "24h Change",
				Value: fmt.Sprintf("%.2f%%", coinList[0].PercentChange24H),
				Short: true,
			},
			slack.AttachmentField{
				Title: coinList[1].ID,
				Value: fmt.Sprintf("$%.2f", coinList[1].PriceUSD),
				Short: true,
			},
			slack.AttachmentField{
				Title: "24h Change",
				Value: fmt.Sprintf("%.2f%%", coinList[1].PercentChange24H),
				Short: true,
			},
			slack.AttachmentField{
				Title: coinList[2].ID,
				Value: fmt.Sprintf("$%.2f", coinList[2].PriceUSD),
				Short: true,
			},
			slack.AttachmentField{
				Title: "24h Change",
				Value: fmt.Sprintf("%.2f%%", coinList[2].PercentChange24H),
				Short: true,
			},
		},
		Color: areWeHappy(coinList[0].PercentChange1H),
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(s.channelID, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
