package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SixofClubsss/dPrediction/prediction"
	"github.com/SixofClubsss/dReams-Bots/common"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	help = `
<a href="https://dreamdapps.io">dReam dApps</a> Bookie

/help - <i>shows this message</i>

<u>Sports commands</u>:
/epl_games - <i>Get Current FIFA games</i>
/nba_games - <i>Get Current NBA games</i>
/nfl_games - <i>Get Current NFL games</i>
/nhl_games - <i>Get Current NHL games</i>
/mma_fights - <i>Get Current MMA fights</i>
/epl_finals - <i>Get FIFA final results</i>
/nba_finals - <i>Get NBA final results</i>
/nfl_finals - <i>Get NFL final results</i>
/nhl_finals - <i>Get NHL final results</i>
/mma_finals - <i>Get MMA final results</i>

<u>Prediction commands</u>:
/btc_usdt - <i>Get BTC-USDT predictions</i>
/dero_usdt - <i>Get DERO-USDT predictions</i>
/xmr_usdt - <i>Get XMR-USDT predictions</i>
/dero_onchain - <i>Get DERO On-Chain predictions</i>

Powered by <a href="http://github.com/civilware/Gnomon">Gnomon</a>`

	bot_name = "Bookie"

	soccer_contract     = "aa57e21c0891a9a99199280284d4a15f2969a0db98166ca2ce8c60a9572e9cba"
	basketball_contract = "ad11377c29a863523c1cc50a33ca13e861cc146a7c0496da58deaa1973e0a39f"
	football_contract   = "f4f89ecf4142145dec38b3e543a10cc1213d13c6d7ca13d01961df93dd2bf3d0"
	hockey_contract     = "c6a7f69ff3f1101a19678b4c28ae5b711c9acc291045049276671493b873dbaa"
	mma_contract        = "faf28fe214271b736f458492295b290b07ae678500f7696419eb02b5969c30b1"

	btcUSDT_contract  = "8a6601afda0d6882bb64bcbbc4848eee474ba9ef0d285469c47415bf4ad469bd"
	deroUSDT_contract = "eaa62b220fa1c411785f43c0c08ec59c761261cb58a0ccedc5b358e5ed2d2c95"
	xmrUSDT_contract  = "db96462400e44fc424c8072b7f328853ed124a8347b7fea8874892a2a58946db"
	onChain_contract  = "a56a89dcbad340b010e028b3b9ff905abaa411c5df60d1ffa8f82f7a9cde6df9"
)

func main() {
	// Initialize logrus to std out
	gnomes.InitLogrusLog(logrus.InfoLevel)

	// Read config.json file
	config := common.ReadConfig(bot_name)

	// Create a new bot instance with API key from config.json file
	bot, err := tgbotapi.NewBotAPI(config.APIKey)
	if err != nil {
		logger.Fatalf("[%s] %s\n", bot_name, err)
	}

	// Parse start flags
	common.Flags()
	logger.Printf("[%s] Authorized on account %s\n", bot_name, bot.Self.UserName)

	// Handle ctrl+c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		common.Close()
		fmt.Println("")
		gnomon.Stop(bot_name)
		logger.Printf("[%s] Exiting\n", bot_name)
		os.Exit(0)
	}()

	// Ping daemon for connection
	rpc.Ping()

	// Start Gnomon
	gnomes.StartGnomon(bot_name, gnomon.DBStorageType(), searchFilters(), 0, 0, nil)

	// Bot update configs
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Start daemon CheckConnection loop
	go common.CheckConnection(bot_name)

	// Watch for updates from telegram
	for update := range updates {
		// If we got a message
		if update.Message != nil && update.Message.Chat.ID == config.ChatID {
			var msg tgbotapi.MessageConfig
			logger.Printf("[%s] [%s] %s\n", bot_name, update.Message.From.UserName, update.Message.Text)
			m := update.Message.Text

			if m == "/epl_games" || m == "/epl_games@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(soccer_contract))

			} else if m == "/nba_games" || m == "/nba_games@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(basketball_contract))

			} else if m == "/nfl_games" || m == "/nfl_games@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(football_contract))

			} else if m == "/nhl_games" || m == "/nhl_games@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(hockey_contract))

			} else if m == "/mma_fights" || m == "/mma_fights@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(mma_contract))

			} else if m == "/epl_finals" || m == "/epl_finals@dReamsBookieBot" {
				finals := prediction.FetchSportsFinal(soccer_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nba_finals" || m == "/nba_finals@dReamsBookieBot" {
				finals := prediction.FetchSportsFinal(basketball_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nfl_finals" || m == "/nfl_finals@dReamsBookieBot" {
				finals := prediction.FetchSportsFinal(football_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nhl_finals" || m == "/nhl_finals@dReamsBookieBot" {
				finals := prediction.FetchSportsFinal(hockey_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/mma_finals" || m == "/mma_finals@dReamsBookieBot" {
				finals := prediction.FetchSportsFinal(mma_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/btc_usdt" || m == "/btc_usdt@dReamsBookieBot" {
				prediction.Predict.Contract.SCID = btcUSDT_contract
				text := prediction.GetPrediction(prediction.Predict.Contract.SCID)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatPrediction(text))

			} else if m == "/dero_usdt" || m == "/dero_usdt@dReamsBookieBot" {
				prediction.Predict.Contract.SCID = deroUSDT_contract
				text := prediction.GetPrediction(prediction.Predict.Contract.SCID)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatPrediction(text))

			} else if m == "/xmr_usdt" || m == "/xmr_usdt@dReamsBookieBot" {
				prediction.Predict.Contract.SCID = xmrUSDT_contract
				text := prediction.GetPrediction(prediction.Predict.Contract.SCID)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatPrediction(text))

			} else if m == "/dero_onchain" || m == "/dero_onchain@dReamsBookieBot" {
				prediction.Predict.Contract.SCID = onChain_contract
				text := prediction.GetPrediction(prediction.Predict.Contract.SCID)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatPrediction(text))

			} else if m == "/help" || m == "/help@dReamsBookieBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, help)
				msg.DisableWebPagePreview = true
			}

			msg.ParseMode = "HTML"
			msg.ReplyToMessageID = update.Message.MessageID
			if msg.Text != "" {
				bot.Send(msg)
			}
		}
	}
}
