package main

import (
	"time"

	"github.com/SixofClubsss/Grokked/grok"
	"github.com/SixofClubsss/dReams-Bots/common"
	"github.com/dReam-dApps/dReams/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	help = `
<a href="https://dreamdapps.io">dReam dApps</a> Grok Bot

<u>Grok Bot Commands</u>:

/help - <i>Shows this message</i>

/grok - <i>Get the current Grok</i>

Powered by <a href="http://github.com/civilware/Gnomon">Gnomon</a>`

	bot_name = "Grok Bot"
)

func main() {
	// Read config.json file
	config := common.ReadConfig(bot_name)

	// Create a new bot instance with API key from config.json file
	bot, err := tgbotapi.NewBotAPI(config.APIKey)
	if err != nil {
		logger.Fatalf("[%s] %s\n", bot_name, err)
	}

	// Grokker service
	go grok.RunGrokker()

	logger.Printf("[%s] Authorized on account %s\n", bot_name, bot.Self.UserName)

	// Bot update configs
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				var msg tgbotapi.MessageConfig
				logger.Printf("[%s] [%s] %s\n", bot_name, update.Message.From.UserName, update.Message.Text)
				if update.Message.Time().Unix() < time.Now().Unix()-30 {
					continue
				}

				switch update.Message.Text {
				case "/grok", "/grok@dReamsGrokBot":
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetGrok(config.SCID))
				case "/help", "/help@dReamsGrokBot":
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, help)
					msg.DisableWebPagePreview = true
				}

				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				if msg.Text != "" {
					bot.Send(msg)
				}
			}
		default:
			if menu.IsClosing() {
				time.Sleep(4 * time.Second)
				logger.Printf("[%s] Closed\n", bot_name)
				return
			}
			time.Sleep(time.Second)
		}
	}
}
