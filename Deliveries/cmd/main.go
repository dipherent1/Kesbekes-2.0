package main

import (
	"fmt"
	config "kesbekes/Config"
	"kesbekes/Infrastructure/bot"
)

func main() {
	fmt.Println("here")
	fmt.Println("here")
	config.EnvInit()
	db := config.ConnectDB()
	fmt.Println(db)
	bot.NewTdLib()
}
