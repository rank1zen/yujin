package rest

import (
	"github.com/gin-gonic/gin"
)


func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	return r
}
