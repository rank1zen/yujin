package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/service"
)

type SummonerHandler struct {
	svc *service.SummonerService
}

func NewSummonerHandler(svc *service.SummonerService) *SummonerHandler {
	return &SummonerHandler{
		svc: svc,
	}
}

func (s *SummonerHandler) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1/summoner")
	{
		v1.GET("/all")
		v1.GET("/all/newest")
		v1.GET("/by_puuid/:puuid", s.byPuuid)
		v1.GET("/by_puuid/:puuid/newest", s.newestByPuuid)

		v1.POST("/renew/:name", s.renew)
	}
}

func (s *SummonerHandler) renew(c *gin.Context) {
	name := c.Param("name")

	msg, err := s.svc.QueryRiot(c, name)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, gin.H{"msg":msg})
}

func (s *SummonerHandler) byPuuid(c *gin.Context) {
	puuid := c.Param("puuid")

	//offset := c.DefaultQuery("offset", "0")
	//limit := c.DefaultQuery("limit", "10")

	summoners, err := s.svc.FindByPuuid(c, puuid)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, summoners)
}

func (s *SummonerHandler) newestByPuuid(c *gin.Context) {
	puuid := c.Param("puuid")

	summoner, err := s.svc.NewestByPuuid(c, puuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, summoner)
}
