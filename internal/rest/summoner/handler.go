package summoner

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Get a list of one summoners' records": "LOL"})
}

func GetRecentRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Get one summoners' most recent record": "LOL"})
}

func PostRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Post one summoner record": "LOL"})
}



