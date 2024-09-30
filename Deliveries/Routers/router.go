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
	Router = gin.Default()

	Bot()
	return Router
}
