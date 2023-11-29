package internal


type RiotRegion string

const (
	RegionNA RiotRegion = "na"
	RegionEUW RiotRegion = "euw"
	RegionKR RiotRegion = "kr"
)

type Summoner struct {
	Region RiotRegion
	PuuId         string
	AccountId     string
	SummonerId    string
	Level         int   
	ProfileIconId int   
	Name          string
	LastRevision  string
}

func (s Summoner) Validate() error {
	// TODO: Implement validation
	return nil
}
