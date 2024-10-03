package routers

import (
	config "kesbekes/Config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Router *gin.Engine
var DB *gorm.DB

func Setuprouter() *gin.Engine {
	db := config.ConnectDB()
	DB = db
	// Migrate the schema
	Migrate()

	Router = gin.Default()

	BotRouter()
	return Router
}
