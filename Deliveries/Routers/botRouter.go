package routers

import (
	config "kesbekes/Config"
	controllers "kesbekes/Deliveries/Controllers"
	"kesbekes/Infrastructure/bot"
	repositories "kesbekes/Repositories"
)

func BotRouter() {
	TgClient := bot.NewTdLib()
	Bot := config.NewBot()

	BotRepo := repositories.NewTelegramRepository(DB)
	botController := controllers.NewBotController(TgClient, Bot, BotRepo)
	botRouter := Router.Group("/bot")

	{
		botRouter.GET("/get10updates", botController.Get10Updates)
	}

	webhookRouter := Router.Group("/webhook")
	{
		webhookRouter.POST("", botController.Webhook)
	}
}
