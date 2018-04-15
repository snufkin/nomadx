package main

import (
	"fmt"
	cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/nlopes/slack"
	"time"
)

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
