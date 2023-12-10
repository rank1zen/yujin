package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal/service"
)

/*
	a := summoner.cast()
	fmt.Printf("%+v\n", a)
*/

func TestTypeCast(t *testing.T) {
	summoner := service.Summoner{
		Region: "na",
		Puuid: "TEST_PUUID",
		AccountId: "TEST_ACCOUNTID",
		SummonerId: "TEST_SUMMONERID",
		Level: 345,
		ProfileIconId: 1033,
		Name: "JOE",
		LastRevision: time.Now(),
		TimeStamp: time.Now(),
	}
	fmt.Printf("%+v\n", summoner)

	a := summoner.Cast()

	fmt.Printf("%+v\n", a)
}
