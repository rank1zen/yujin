package service



type SummonerRepo interface {
	Create() error
	Find() ()
	FindR()

}

type Summoner struct {
	repo SummonerRepo
}

func NewSummoner(repo SummonerRepo) *Summoner {
	return &Summoner{
		repo: repo,
	}
}


