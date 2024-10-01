package routers

import (
	controllers "kesbekes/Deliveries/Controllers"
	"kesbekes/Infrastructure/bot"
)

func Bot() {
	botClient := bot.NewTdLib()
	botController := controllers.NewBotController(botClient)
	botRouter := Router.Group("/bot")
	{
		botRouter.GET("/get10updates", botController.Get10Updates)
	}

	webhookRouter := Router.Group("/webhook")
	{
		webhookRouter.POST("", botController.Webhook)
	}
}
