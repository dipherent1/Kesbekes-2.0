package config

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewBot() *tgbotapi.BotAPI {

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}

	wh, _ := tgbotapi.NewWebhook("https://f2f4-196-190-62-133.ngrok-free.app/webhook")

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram webhook error: %s", info.LastErrorMessage)
	}
	return bot
}
