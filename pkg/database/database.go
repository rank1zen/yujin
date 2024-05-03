package database

import (
	"context"
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	soloQueueType = 420
	soloqOption   = lol.MatchListOptions{Queue: &soloQueueType}
)

// This is a wrapper for exclusivly "QUERY" logic
type pgxDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}

// IRecordQuery represents methods for a specific table
type IRecordQuery[M any] interface {
	GetRecords(ctx context.Context, filters ...RecordFilter) ([]*M, error)
	CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error)
	InsertRecords(ctx context.Context, records []M) (int64, error)
	DeleteRecords(ctx context.Context) error
}

type DB interface {
	SummonerV4() IRecordQuery[SummonerRecord]
	LeagueV4() IRecordQuery[LeagueRecord]
	MatchV5() IRecordQuery[MatchRecord]
	MatchV5Ban() IRecordQuery[MatchBanRecord]
	MatchV5Team() IRecordQuery[MatchTeamRecord]
	MatchV5Objective() IRecordQuery[MatchObjectiveRecord]
	MatchV5Paricipant() IRecordQuery[MatchParticipantRecord]

	FetchAndInsertSummoner(ctx context.Context, puuid string) error
	FetchAndInsertRank(ctx context.Context, summmonerId string) error
	FetchAndInsertMatches(ctx context.Context, puuid string) error
	FetchAndInsertAllMatches(ctx context.Context, puuid string) error

	Close()
}

// db represents a collection of things to access database things
type db struct {
	summonerV4        *summonerV4Query
	leagueV4          *leagueV4Query
	matchV5           *matchV5Query
	matchV5Ban        *matchV5BanQuery // BRUH
	matchV5Obj        *matchV5ObjQuery // BRUH
	matchV5Team       *matchV5TeamQuery
	matchV5Paricipant *matchV5ParticipantQuery

	pgx pgxDB
	gc  *golio.Client
}

func NewDB(ctx context.Context, url string, gc *golio.Client) (DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &db{
		summonerV4:        &summonerV4Query{db: pool},
		leagueV4:          nil,
		matchV5:           &matchV5Query{db: pool},
		matchV5Ban:        nil,
		matchV5Obj:        nil,
		matchV5Team:       nil,
		matchV5Paricipant: nil,
		pgx:               pool,
		gc:                gc,
	}, nil
}

func (d *db) SummonerV4() IRecordQuery[SummonerRecord] {
	return d.summonerV4
}

func (d *db) LeagueV4() IRecordQuery[LeagueRecord] {
	return d.leagueV4
}

func (d *db) MatchV5() IRecordQuery[MatchRecord] {
	return d.matchV5
}

func (d *db) MatchV5Ban() IRecordQuery[MatchBanRecord] {
	return d.matchV5Ban
}

func (d *db) MatchV5Team() IRecordQuery[MatchTeamRecord] {
	return d.matchV5Team
}

func (d *db) MatchV5Paricipant() IRecordQuery[MatchParticipantRecord] {
	return d.matchV5Paricipant
}

func (d *db) MatchV5Objective() IRecordQuery[MatchObjectiveRecord] {
	return d.matchV5Obj
}

func (d *db) FetchAndInsertSummoner(ctx context.Context, puuid string) error {
	summ, err := d.gc.Riot.LoL.Summoner.GetByPUUID(puuid)
	if err != nil {
		return fmt.Errorf("GetByPUUID: %w", err)
	}

	record, err := newSummonerRecord(summ)
	if err != nil {
		return fmt.Errorf("newSummonerRecord: %w", err)
	}

	_, err = d.summonerV4.InsertRecords(ctx, record)
	if err != nil {
		return fmt.Errorf("insert records: %w", err)
	}

	return nil
}

func (d *db) FetchAndInsertRank(ctx context.Context, summonerId string) error {
	return fmt.Errorf("not implemented")
}

func (d *db) FetchAndInsertMatches(ctx context.Context, puuid string) error {
	var m MatchThingy

	ids, err := d.gc.Riot.LoL.Match.List(puuid, 0, 5, &soloqOption)
	if err != nil {
		return fmt.Errorf("fetch match ids: %w", err)
	}

	err = d.fetchMatch(ctx, ids, &m)
	if err != nil {
		return fmt.Errorf("fetch matches: %w", err)
	}

	err = d.insertMatches(ctx, &m)
	if err != nil {
		return fmt.Errorf("insert matches: %w", err)
	}

	return nil
}

func (d *db) FetchAndInsertAllMatches(ctx context.Context, puuid string) error {
	return fmt.Errorf("not implemented")
	var m MatchThingy

	var ids []string
	ch := d.gc.Riot.LoL.Match.ListStream(puuid)
	for match := range ch {
		if match.Error != nil {
			return fmt.Errorf("okoko :%w", match.Error)
		}

		ids = append(ids, match.MatchID)
	}

	err := d.fetchMatch(ctx, ids, &m)
	if err != nil {
		return fmt.Errorf("")
	}

	return nil
}

// Close closes the DB connection
func (d *db) Close() {
	d.pgx.Close()
}

func (d *db) fetchMatch(ctx context.Context, matchIds []string, m *MatchThingy) error {
	for _, id := range matchIds {
		match, err := d.gc.Riot.LoL.Match.Get(id)
		if err != nil {
			return fmt.Errorf("fetch match :%w", err)
		}

		m.info = append(m.info, NewMatchRecord(match.Info)...)
		m.team = append(m.team, NewMatchTeamRecord(match.Info.Participants...)...)
	}

	return nil
}

func (d *db) insertMatches(ctx context.Context, m *MatchThingy) error {
	// prob want to do this in a transaction or something
	// also like we may not want to put all our eggs in one basket
	_, err := d.matchV5.InsertRecords(ctx, m.info)
	if err != nil {
		return nil
	}

	return nil
}

type MatchThingy struct {
	info []MatchRecord
	team []MatchTeamRecord
}
