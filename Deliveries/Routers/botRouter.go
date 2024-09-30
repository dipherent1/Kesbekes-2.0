package routers

import (
	controllers "kesbekes/Deliveries/Controllers"
	"kesbekes/Infrastructure/bot"
)

func Bot() {
	botRouter := Router.Group("/bot")
	{
		botClient := bot.NewTdLib()
		botController := controllers.NewBotController(botClient)

		botRouter.GET("/get10updates", botController.Get10Updates)
	}
}
