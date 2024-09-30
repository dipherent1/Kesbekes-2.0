package controllers

import (
	"kesbekes/Infrastructure/bot"

	"github.com/gin-gonic/gin"
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
