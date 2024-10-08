package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/pgxutil"
	"github.com/rank1zen/yujin/internal/riot"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	pool   *pgxpool.Pool
	riot   *riot.Client
	tracer trace.Tracer
}

func NewDB(ctx context.Context, url string) (*DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pgxCfg.ConnConfig.Tracer = &tracer{}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	riot := &riot.Client{}

	return &DB{
		pool: pool,
		riot: riot,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func riotGetName(ctx context.Context, riot *riot.Client, puuid string) (string, error) {
	account, err := riot.AccountGetByPuuid(ctx, puuid)
	if err != nil {
		return "", err
	}

	return account.GameName + "#" + account.TagLine, nil
}

func riotGetSummonerID(ctx context.Context, riot *riot.Client, puuid string) (string, error) {
	account, err := riot.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return "", err
	}

	return account.Id, nil
}

func riotUnixToDate(ts int64) time.Time {
	return time.Unix(0, ts)
}

func riotDurationToInterval(dur int) time.Duration {
	return time.Duration(dur) * time.Second
}

func dbGetItemIconUrls(ctx context.Context, db pgxutil.Conn, ids [7]int) [7]*string {
	var urls [7]*string
	err := db.QueryRow(ctx, `SELECT get_item_icon_urls($1)`, ids).Scan(&urls)
	if err != nil {
		logging.FromContext(ctx).Sugar().DPanic(err)
	}
	return urls
}

func dbGetChampionIconUrl(ctx context.Context, db pgxutil.Conn, id int) string {
	var url string
	err := db.QueryRow(ctx, `SELECT get_item_icon_urls($1)`, id).Scan(&url)
	if err != nil {
		logging.FromContext(ctx).Sugar().DPanic(err)
	}
	return url
}

func dbGetSummonersIconUrls(ctx context.Context, db pgxutil.Conn, ids [2]int) [2]string {
	var urls [2]string
	err := db.QueryRow(ctx, `SELECT get_summoners_icon_urls($1)`, ids).Scan(&urls)
	if err != nil {
		logging.FromContext(ctx).Sugar().DPanic(err)
	}
	return urls
}

func dbGetRuneIconUrl(ctx context.Context, db pgxutil.Conn, id int) string {
	var url string
	err := db.QueryRow(ctx, `SELECT get_rune_icon_urls($1)`, id).Scan(&url)
	if err != nil {
		logging.FromContext(ctx).Sugar().DPanic(err)
	}
	return url
}

func dbGetRuneTreeIconUrl(ctx context.Context, db pgxutil.Conn, id int) string {
	var url string
	err := db.QueryRow(ctx, `SELECT get_rune_tree_icon_urls($1)`, id).Scan(&url)
	if err != nil {
		logging.FromContext(ctx).Sugar().DPanic(err)
	}
	return url
}

type Account struct {
	Puuid      string
	SummonerId string
	Name       string
	TagLine    string
}

func (db *DB) GetAccount(ctx context.Context, name string) (*Account, error) {
	getAccount := func(gamename, tagline string) (*Account, error) {
		var ids Account
		err := db.pool.QueryRow(ctx, `
		SELECT
			summoner_id,
			puuid,
			name,
			tagline
		FROM
			riot_accounts
		WHERE 1=1
			AND UPPER(name) = UPPER($1)
			AND UPPER(tagline) = UPPER($2);
		`, gamename, tagline).Scan(&ids.SummonerId, &ids.Puuid, &ids.Name, &ids.TagLine)
		if err != nil {
			return nil, fmt.Errorf("getting db: %w", err)
		}
		return &ids, nil
	}

	parts := strings.SplitN(name, "-", 2)
	gamename, tagline := parts[0], parts[1]
	var found bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM
			riot_accounts
		WHERE 1=1
			AND UPPER(name) = UPPER($1)
			AND UPPER(tagline) = UPPER($2)
	);
	`, gamename, tagline).Scan(&found)
	if err != nil {
		return nil, fmt.Errorf("checking db: %w", err)
	}

	if found {
		return getAccount(gamename, tagline)
	}

	acc, err := db.riot.AccountGetByRiotId(ctx, gamename, tagline)
	if err != nil {
		return nil, fmt.Errorf("fetching account: %w", err)
	}

	summ, err := db.riot.GetSummonerByPuuid(ctx, acc.Puuid)
	if err != nil {
		return nil, fmt.Errorf("fetching summoner: %w", err)
	}

	vals := map[string]any{
		"puuid":       summ.Puuid,
		"summoner_id": summ.Id,
		"name":        acc.GameName,
		"tagline":     acc.TagLine,
	}

	err = pgxutil.QueryInsertRow(ctx, db.pool, "riot_accounts", vals)
	if err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	return getAccount(gamename, tagline)
}
