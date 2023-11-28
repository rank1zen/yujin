package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/rest/summoner"
)


func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/summoner/:region/:puuid", summoner.GetRecords)
		v1.GET("/summoner/:region/:puuid/recent", summoner.GetRecentRecord)
		v1.POST("/summoner/:region/", summoner.PostRecord)
	}

	return r
}
