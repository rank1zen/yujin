package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	_ = gin.Default()


	_ = "postgres://gordon:kop123456@localhost:5432"
}
