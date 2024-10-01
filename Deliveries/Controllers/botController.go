package controllers

import (
	"kesbekes/Infrastructure/bot"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotController struct {
	TdlibClient *bot.TdLib
}

func NewBotController(tdlibClient *bot.TdLib) *BotController {
	return &BotController{
		TdlibClient: tdlibClient,
	}
}

func (b *BotController) Get10Updates(c *gin.Context) {
	b.TdlibClient.Get10Updates()
}

func (b *BotController) Webhook(c *gin.Context) {
	var update tgbotapi.Update

	// Parse the incoming request body (JSON) into the Update struct
	if err := c.ShouldBindJSON(&update); err != nil {
		log.Printf("Error decoding incoming update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Handle the forwarded message
	if update.Message != nil && update.Message.ForwardFromChat != nil {
		chatID := update.Message.ForwardFromChat.ID
		log.Printf("Received forwarded message from chat ID: %d", chatID)
	}

	// Respond back to Telegram
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
