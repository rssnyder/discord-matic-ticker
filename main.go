package main

import (
	"flag"
	"os"
	"sync"

	env "github.com/caitlinelfring/go-env-default"
	log "github.com/sirupsen/logrus"
)

var (
	logger  = log.New()
	address *string
)

func init() {
	// initialize logging
	logLevel := flag.Int("logLevel", 0, "defines the log level. 0=production builds. 1=dev builds.")
	address = flag.String("address", "localhost:8080", "address:port to bind http server to.")
	flag.Parse()
	logger.Out = os.Stdout
	switch *logLevel {
	case 0:
		logger.SetLevel(log.InfoLevel)
	default:
		logger.SetLevel(log.DebugLevel)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	m := NewManager(*address)

	// check for inital bots
	if os.Getenv("DISCORD_BOT_TOKEN") != "" {
		s := addInitialToken()
		m.addToken(s)
	}

	// wait forever
	wg.Wait()
}

func addInitialToken() *Token {
	var matic *Token

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		logger.Fatal("Discord bot token is not set! Shutting down.")
	}

	contract := os.Getenv("CONTRACT")
	name := os.Getenv("NAME")
	nickname := env.GetBoolDefault("SET_NICKNAME", false)
	frequency := env.GetIntDefault("FREQUENCY", 60)
	currency := env.GetDefault("CURRENCY", "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174")

	matic = NewToken(contract, token, name, nickname, frequency, currency)

	return matic
}
