package controllers

import (
	domains "kesbekes/Domains"
	"kesbekes/Infrastructure/bot"
	repositories "kesbekes/Repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotController struct {
	TdlibClient *bot.TdLib
	Bot         *tgbotapi.BotAPI
	BotRepo     *repositories.TelegramRepository
}

func NewBotController(tdlibClient *bot.TdLib, bot *tgbotapi.BotAPI, BotRepo *repositories.TelegramRepository) *BotController {
	return &BotController{
		TdlibClient: tdlibClient,
		Bot:         bot,
		BotRepo:     BotRepo,
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
		// Check if the chat is already in the database
		chatinfo := &domains.ChatInfo{
			Name:     update.Message.ForwardFromChat.Title,
			Username: update.Message.ForwardFromChat.UserName,
			ChatID:   update.Message.ForwardFromChat.ID,
		}

		// Save the chat to the database
		err := b.BotRepo.StoreChat(chatinfo)
		if err != nil {
			log.Printf("Error saving chat: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error saving chat"))
			return
		}

		// Respond back to Telegram
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Chat saved"))
		return
	}

	// Respond back to Telegram
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
