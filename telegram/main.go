package main

import (
	"fmt"
	"log"

	"github.com/SixofClubsss/dReams/prediction"
	"github.com/SixofClubsss/dReams/rpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	help = `
<a href="https://dreamtables.net">dReam Tables</a> Bookie

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

	soccer_contract     = "aa57e21c0891a9a99199280284d4a15f2969a0db98166ca2ce8c60a9572e9cba"
	basketball_contract = "ad11377c29a863523c1cc50a33ca13e861cc146a7c0496da58deaa1973e0a39f"
	football_contract   = "f4f89ecf4142145dec38b3e543a10cc1213d13c6d7ca13d01961df93dd2bf3d0"
	hockey_contract     = "c6a7f69ff3f1101a19678b4c28ae5b711c9acc291045049276671493b873dbaa"
	mma_contract        = "faf28fe214271b736f458492295b290b07ae678500f7696419eb02b5969c30b1"

	btcUSDT_contract  = "c89c2f514300413fd6922c28591196a7c48b42b07e7f4d7d8d9f7643e253a6ff"
	deroUSDT_contract = "eaa62b220fa1c411785f43c0c08ec59c761261cb58a0ccedc5b358e5ed2d2c95"
	xmrUSDT_contract  = "db96462400e44fc424c8072b7f328853ed124a8347b7fea8874892a2a58946db"
	onChain_contract  = "a56a89dcbad340b010e028b3b9ff905abaa411c5df60d1ffa8f82f7a9cde6df9"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(readConfig())
	if err != nil {
		log.Panic("[Telegram-Bot]", err)
	}

	bot.Debug = flags()
	log.Printf("[Telegram-Bot] Authorized on account %s", bot.Self.UserName)

	rpc.Ping()
	startGnomon()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go checkConnection()

	for update := range updates {
		if update.Message != nil { // If we got a message
			var msg tgbotapi.MessageConfig
			log.Printf("[Telegram-Bot] [%s] %s", update.Message.From.UserName, update.Message.Text)
			m := update.Message.Text

			if m == "/epl_games" || m == "/epl_games@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(rpc.Signal.Daemon, soccer_contract))

			} else if m == "/nba_games" || m == "/nba_games@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(rpc.Signal.Daemon, basketball_contract))

			} else if m == "/nfl_games" || m == "/nfl_games@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(rpc.Signal.Daemon, football_contract))

			} else if m == "/nhl_games" || m == "/nhl_games@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(rpc.Signal.Daemon, hockey_contract))

			} else if m == "/mma_fights" || m == "/mma_fights@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetBook(rpc.Signal.Daemon, mma_contract))

			} else if m == "/epl_finals" || m == "/epl_finals@dReamTables_bot" {
				finals, _ := rpc.FetchSportsFinal(rpc.Signal.Daemon, soccer_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nba_finals" || m == "/nba_finals@dReamTables_bot" {
				finals, _ := rpc.FetchSportsFinal(rpc.Signal.Daemon, basketball_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nfl_finals" || m == "/nfl_finals@dReamTables_bot" {
				finals, _ := rpc.FetchSportsFinal(rpc.Signal.Daemon, football_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/nhl_finals" || m == "/nhl_finals@dReamTables_bot" {
				finals, _ := rpc.FetchSportsFinal(rpc.Signal.Daemon, hockey_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/mma_finals" || m == "/mma_finals@dReamTables_bot" {
				finals, _ := rpc.FetchSportsFinal(rpc.Signal.Daemon, mma_contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, formatFinals(finals))

			} else if m == "/btc_usdt" || m == "/btc@dReamTables_bot" {
				prediction.PredictControl.Contract = btcUSDT_contract
				text := prediction.GetPrediction(rpc.Signal.Daemon, prediction.PredictControl.Contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, HtmlTextParse(text))

			} else if m == "/dero_usdt" || m == "/dero@dReamTables_bot" {
				prediction.PredictControl.Contract = deroUSDT_contract
				text := prediction.GetPrediction(rpc.Signal.Daemon, prediction.PredictControl.Contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, HtmlTextParse(text))

			} else if m == "/xmr_usdt" || m == "/xmr@dReamTables_bot" {
				prediction.PredictControl.Contract = xmrUSDT_contract
				text := prediction.GetPrediction(rpc.Signal.Daemon, prediction.PredictControl.Contract)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, HtmlTextParse(text))

			} else if m == "/dero_onchain" || m == "/dero_onchain@dReamTables_bot" {
				prediction.PredictControl.Contract = onChain_contract
				text := prediction.GetPrediction(rpc.Signal.Daemon, prediction.PredictControl.Contract)
				fmt.Println(text)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, HtmlTextParse(text))

			} else if m == "/help" || m == "/help@dReamTables_bot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, help)
				msg.DisableWebPagePreview = true
			}

			msg.ParseMode = "HTML"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
