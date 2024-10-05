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

	wh, _ := tgbotapi.NewWebhook("https://e850-196-190-60-154.ngrok-free.app/webhook")

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
