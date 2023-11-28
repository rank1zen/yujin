package internal

type MSummoner struct {
	PuuId         string
	AccountId     string
	SummonerId    string
	Level         int   
	ProfileIconId int   
	Name          string
}

func (s MSummoner) Validate() error {
	// TODO: Implement validation
	return nil
}
