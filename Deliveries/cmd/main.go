package main

import (
	"fmt"
	config "kesbekes/Config"
)

func main() {
	config.EnvInit()
	db := config.ConnectDB()
	fmt.Println(db)
}
