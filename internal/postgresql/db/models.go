package db

type Region string

const (
	RegionNA Region = "na"
	RegionEUW Region = "euw"
	RegionKR Region = "kr"
)

type Summoner struct {
	PuuId         string `json:"puuid" gorm:"primaryKey"`
	AccountId     string `json:"accountId"`
	SummonerId    string `json:"summonerId"`
	Level         int    `json:"summonerLevel"`
	ProfileIconId int    `json:"profileIconId"`
	Name          string `json:"name"`
}
