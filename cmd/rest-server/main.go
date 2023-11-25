package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/rest"
)

func main() {
	r := gin.New()

	r.GET("/", rest.HomeHandler)

	r.Run()
}

