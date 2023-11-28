package internal

type CreateSummonerParams struct {
	Region string
	Puuid string
	Limit int32
	Offset int32
}

type FindSummonerParams struct {
	Region string
	Puuid string
	Limit int32
	Offset int32
}
