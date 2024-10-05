package controllers

import (
	"fmt"
	domains "kesbekes/Domains"
	"kesbekes/Infrastructure/bot"
	repositories "kesbekes/Repositories"
	"log"
	"net/http"
	"strings"

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

var isPreferenses = false
var isDeletePreferenses = false
var isDeleteChat = false

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
			Username: "@" + update.Message.ForwardFromChat.UserName,
			ChatID:   update.Message.ForwardFromChat.ID,
		}

		// Save the chat to the database
		err := b.BotRepo.StoreChat(chatinfo, int64(update.Message.From.ID))
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

	// Handle the message start command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "start" {
		// Check if the user is already in the database
		user := &domains.User{
			Username: update.Message.From.UserName,
			FistName: update.Message.From.FirstName,
			LastName: update.Message.From.LastName,
			Role:     "user",
			UserID:   int64(update.Message.From.ID),
		}
		err := b.BotRepo.StoreUser(user)
		if err != nil {
			log.Printf("Error saving user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error saving user"))
			return
		}

		// Respond back to Telegram
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "User saved"))
		return
	}

	// Handle the message add perfereces command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "addpreferenses" {
		isPreferenses = true
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Enter your preferenses in the following format:\n preferense1\n preferense2\n preferense3"))
		return
	}

	if isPreferenses {
		isPreferenses = false
		preferenses := []string{}
		// Split the message by new line and add each line to the preferenses slice
		lines := strings.Split(update.Message.Text, "\n")
		for _, line := range lines {
			preferenses = append(preferenses, line)
		}

		user := &domains.User{
			UserID:      update.Message.From.ID,
			Preferenses: preferenses,
		}
		err := b.BotRepo.UpdateUserPreferences(user)
		if err != nil {
			log.Printf("Error saving preferenses: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error saving preferenses"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Preferenses saved"))
		return
	}

	// Handle the message get preferenses command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "getpreferenses" {
		userID := int64(update.Message.From.ID)
		preferenses, err := b.BotRepo.GetUserPreferences(userID)
		if err != nil {
			log.Printf("Error getting preferenses: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting preferenses"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(preferenses, "\n")))
		return
	}

	// delelet preferennse command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "deletepreferenses" {
		isDeletePreferenses = true
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Enter the preferense you want to delete one by one"))
		return
	}

	if isDeletePreferenses {
		isDeletePreferenses = false
		deleteWord := update.Message.Text
		userID := int64(update.Message.From.ID)
		err := b.BotRepo.DeleteUserPreferences(userID, deleteWord)
		if err != nil {
			log.Printf("Error deleting preferenses: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error deleting preferenses"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Preferense deleted"))
		return
	}

	//get my chats command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "getmychats" {
		userID := int64(update.Message.From.ID)
		chats, err := b.BotRepo.GetUserChats(userID)
		if err != nil {
			log.Printf("Error getting chats: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting chats"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Your chats are:"))
		for _, chat := range chats {
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, chat.Username))
		}
		return
	}

	//delete chat command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "deletechat" {
		isDeleteChat = true
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Enter the chat username you want to delete"))
		return
	}

	if isDeleteChat {
		isDeleteChat = false
		deleteWord := update.Message.Text
		userID := int64(update.Message.From.ID)
		err := b.BotRepo.DeleteUserChat(userID, deleteWord)
		if err != nil {
			log.Printf("Error deleting chat: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error deleting chat"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Chat deleted"))
		return
	}

	// handel message listen command
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "listen" {
		// get chats
		userID := int64(update.Message.From.ID)
		chats, err := b.BotRepo.GetUserChats(userID)
		if err != nil {
			log.Printf("Error getting chats: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			//send using telegram
			b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting chats"))
			return
		}

		// get chat id in an array
		chatIDs := []int64{}
		for _, chat := range chats {
			chatIDs = append(chatIDs, chat.ChatID)
		}

		fmt.Println(chatIDs)

		// Run the listener in a separate goroutine so it doesn't block the main process
		go b.TdlibClient.Listen(chatIDs, update.Message.Chat.ID)

		// Acknowledge that the bot started listening
		b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Now listening to your chats."))
	}

	// Respond back to Telegram
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
