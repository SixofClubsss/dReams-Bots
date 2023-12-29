package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SixofClubsss/Grokked/grok"
	"github.com/SixofClubsss/dReams-Bots/common"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
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

	// Get Grokked SCID code for Gnomon filter, fatal if not found
	filter := rpc.GetSCCode(grok.GROKSCID)
	if filter == "" {
		logger.Fatalf("[%s] Could not get Gnomon filter\n", bot_name)
	}

	// Start Gnomon
	gnomes.StartGnomon(bot_name, gnomon.DBStorageType(), []string{filter}, 0, 0, nil)

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

			if m == "/grok" || m == "/grok@dReamsGrokBot" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetGrok(config.SCID))
			} else if m == "/help" || m == "/help@dReamsGrokBot" {
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
