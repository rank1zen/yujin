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
		v1.GET("/:region/all")
		v1.GET("/:region/all/newest")
		v1.GET("/:region/by_puuid/:puuid", s.byPuuid)
		v1.GET("/:region/by_puuid/:puuid/newest", s.newestByPuuid)
	}
}

func (s *SummonerHandler) byPuuid(c *gin.Context) {
	var req ReqByPuuid
	if err := c.ShouldBindUri(&req); err != nil {
		return
	}

	//offset := c.DefaultQuery("offset", "0")
	//limit := c.DefaultQuery("limit", "10")

	summoners, err := s.svc.FindByPuuid(c, req.Puuid)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, summoners)
}

func (s *SummonerHandler) newestByPuuid(c *gin.Context) {
	var req ReqByPuuid
	c.ShouldBindUri(&req)

	summoner, err := s.svc.NewestByPuuid(c, req.Puuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, summoner)
}

type ReqByPuuid struct {
	Region string
	Puuid  string
}
