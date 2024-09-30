package main

import (
	config "kesbekes/Config"
	routers "kesbekes/Deliveries/Routers"
)

func main() {

	config.EnvInit()
	r := routers.Setuprouter()
	r.Run(":8080")

}
