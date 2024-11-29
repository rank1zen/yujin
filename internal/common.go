package internal

type Season int

const (
	Season2020 Season = 1
	SeasonAll  Season = -1
)

type Items [7]*int

type Summoners [2]int

type Champion int

func GetBannedChampion(championID int) *Champion {
	if championID == 0 {
		return nil
	} else {
		id := championID
		return (*Champion)(&id)
	}
}

type Runes [11]int

