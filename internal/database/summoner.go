package database

import (
	"context"

	"github.com/rank1zen/yujin/internal/ddragon"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

type SummonerChampion struct {
	ChampionIcon      string
	Spell1Icon        string
	Spell2Icon        string
	RunePrimaryIcon   string
	RuneSecondaryIcon string
}

func summonerGetChampion(ctx context.Context, db pgxutil.Conn, dst *SummonerChampion, matchID, puuid string) error {
	var champion, spell1, spell2, runePrimary, runeSecondary int
	err := db.QueryRow(ctx, `
	SELECT
		FORMAT('https://cdn.communitydragon.org/14.16.1/champion/%s/square', champion_id) AS champion_icon_url,
		array[item0_id, item1_id, item2_id, item3_id, item4_id, item5_id] as items,
		array[spell1_id, spell2_id] as spells,
		rune_primary_keystone AS rune_primary,
		rune_secondary_path AS rune_secondary
	FROM
		profile_matches
	WHERE 1=1
		AND match_id = $1;
		AND puuid = $1;
	`, matchID, puuid).Scan(&champion, &spell1, &spell2, &runePrimary, &runeSecondary)
	if err != nil {
		return err
	}

	dst.ChampionIcon = ddragon.GetChampionIconUrl(champion)
	dst.Spell1Icon = ddragon.GetSummonerSpellUrl(spell1)
	dst.Spell2Icon = ddragon.GetSummonerSpellUrl(spell2)
	dst.RunePrimaryIcon = ddragon.GetRuneIconUrl(runePrimary)
	dst.RuneSecondaryIcon = ddragon.GetRuneTreeIconUrl(runeSecondary)
	return nil
}
