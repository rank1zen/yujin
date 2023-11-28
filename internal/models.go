package internal

type Summoner struct {
	PuuId         string
	AccountId     string
	SummonerId    string
	Level         int   
	ProfileIconId int   
	Name          string
}

func (s Summoner) Validate() error {
	// TODO: Implement validation
	return nil
}
