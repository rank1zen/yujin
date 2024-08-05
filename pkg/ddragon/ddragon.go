package ddragon

import (
	"fmt"

	"github.com/a-h/templ"
)

// NOTE: static data and images are hard af

func GetItemIconUrl(itemId int) templ.SafeURL {
	u := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/item/%d.png", itemId)
	return templ.URL(u)
}

func GetSummonerProfileIconUrl(iconId int32) templ.SafeURL {
	u := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/profileicon/%d.png", iconId)
	return templ.URL(u)
}

func GetChampionIconUrl(championName string) templ.SafeURL {
	u := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/%s.png", championName)
	return templ.URL(u)
}

func GetChampionSpellUrl() templ.SafeURL {
	return templ.URL("https://www.svgrepo.com/show/379925/alert-error.svg")
}

func GetSummonerSpellUrl(spellId int) templ.SafeURL {
	return templ.URL("https://www.svgrepo.com/show/379925/alert-error.svg")
}

func GetRuneIconUrl() templ.SafeURL {
	return templ.URL("https://www.svgrepo.com/show/379925/alert-error.svg")
}
