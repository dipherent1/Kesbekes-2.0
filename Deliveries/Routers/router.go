package routers

import (
	config "kesbekes/Config"
	ai "kesbekes/Infrastructure/AI"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var Router *gin.Engine
var DB *gorm.DB
var Bot *tgbotapi.BotAPI
var AI *ai.AI

func Setuprouter() *gin.Engine {
	db := config.ConnectDB()
	DB = db
	AI = ai.NewAI()
	// Migrate the schema
	Migrate()

	Router = gin.Default()

	BotRouter()
	return Router
}
