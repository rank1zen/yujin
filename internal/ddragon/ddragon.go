package ddragon

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/rank1zen/yujin/internal/logging"
)

// NOTE: static data and images are hard af

//go:embed assets/*.json
var assets embed.FS

type SummonerSpell struct {
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
	Id       int    `json:"id"`
}

type Perk struct {
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
	Id       int    `json:"id"`
}

var (
	IdToSpell = map[int]SummonerSpell{}
	IdToPerk  = map[int]Perk{}
)

func init() {
	f, _ := fs.Sub(assets, "assets")
	spells, err := f.Open("summoner-spells.json")
	if err != nil {
		logging.Get().Sugar().DPanicf("opening summoner spells: %v", err)
	}

	perks, err := f.Open("perks.json")
	if err != nil {
		logging.Get().Sugar().DPanicf("opening summoner spells: %v", err)
	}

	var m []SummonerSpell
	err = json.NewDecoder(spells).Decode(&m)
	if err != nil {
		logging.Get().Sugar().DPanicf("decoding: %v", err)
	}


	var p []Perk
	err = json.NewDecoder(perks).Decode(&p)
	if err != nil {
		logging.Get().Sugar().DPanicf("decoding: %v", err)
	}

	for _, spell := range m {
		IdToSpell[spell.Id] = spell
	}

	for _, perk := range p {
		IdToPerk[perk.Id] = perk
	}
}

func GetItemIconUrl(itemId int) string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/item/%d.png", itemId)
}

func GetSummonerProfileIconUrl(iconId int) string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/profileicon/%d.png", iconId)
}

func GetChampionIconUrl(championName string) string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/%s.png", championName)
}

func GetSummonerSpellUrl(spellId int) string {
	path, _ := strings.CutPrefix(IdToSpell[spellId].IconPath, "/lol-game-data/assets")
	return "https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default" + strings.ToLower(path)
}

func GetRuneIconUrl(runeId int) string {
	path, _ := strings.CutPrefix(IdToPerk[runeId].IconPath, "/lol-game-data/assets")
	return "https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default" + strings.ToLower(path)
}
