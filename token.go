package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/rssnyder/discord-matic-ticker/utils"
)

type Token struct {
	Contract  string        `json:"contract"`  // stock symbol
	Name      string        `json:"name"`      // override for symbol as shown on the bot
	Nickname  bool          `json:"nickname"`  // flag for changing nickname
	Frequency time.Duration `json:"frequency"` // how often to update in seconds
	Currency  string        `json:"currency"`  // how often to update in seconds
	Price     int           `json:"-"`
	token     string        `json:"-"` // discord token
	close     chan int      `json:"-"`
}

// NewToken saves information about the stock and starts up a watcher on it
func NewToken(contract string, token string, name string, nickname bool, frequency int, currency string) *Token {
	s := &Token{
		Contract:  contract,
		Name:      name,
		Nickname:  nickname,
		Frequency: time.Duration(frequency) * time.Second,
		Currency:  currency,
		token:     token,
		close:     make(chan int, 1),
	}

	// spin off go routine to watch the price
	go s.watchTokenPrice()
	return s
}

// Shutdown sends a signal to shut off the goroutine
func (s *Token) Shutdown() {
	s.close <- 1
}

func (s *Token) watchTokenPrice() {

	// create a new discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + s.token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// show as online
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening discord connection,", err)
		return
	}

	// Get guides for bot
	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println("Error getting guilds: ", err)
		s.Nickname = false
	}

	logger.Infof("Watching token price for %s", s.Name)
	ticker := time.NewTicker(s.Frequency)

	// continuously watch
	for {
		select {
		case <-s.close:
			logger.Infof("Shutting down price watching for %s", s.Name)
			return
		case <-ticker.C:
			logger.Infof("Fetching stock price for %s", s.Name)

			// save the price struct & do something with it
			priceData, err := utils.GetTokenPrice(s.Contract, s.Currency)
			if err != nil {
				logger.Errorf("Unable to fetch stock price for %s", s.Name)
			}

			var fmtPriceRaw float64

			if fmtPriceRaw, err = strconv.ParseFloat(priceData.Totokenamount, 64); err != nil {
				logger.Errorf("Error with price format for %s", s.Name)
			}

			fmtPrice := fmtPriceRaw / 10000000

			if s.Nickname {
				// update nickname instead of activity
				var nickname string
				var activity string

				// format nickname & activity
				nickname = fmt.Sprintf("%s - $%.2f", s.Name, fmtPrice)
				activity = "Using USDC on 1inch"

				// Update nickname in guilds
				for _, g := range guilds {
					err = dg.GuildMemberNickname(g.ID, "@me", nickname)
					if err != nil {
						fmt.Println("Error updating nickname: ", err)
						continue
					}
					logger.Infof("Set nickname in %s: %s", g.Name, nickname)
				}

				err = dg.UpdateGameStatus(0, activity)
				if err != nil {
					logger.Error("Unable to set activity: ", err)
				} else {
					logger.Infof("Set activity: %s", activity)
				}

			} else {
				activity := fmt.Sprintf("%s - $%.2f", s.Name, fmtPrice)

				err = dg.UpdateGameStatus(0, activity)
				if err != nil {
					logger.Error("Unable to set activity: ", err)
				} else {
					logger.Infof("Set activity: %s", activity)
				}

			}

		}

	}

}
