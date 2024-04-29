package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlQueryBuilder(t *testing.T) {
        for _, test := range []struct{
                filters []RecordFilter
                want string
        }{
                {
                        []RecordFilter{
                                {
                                        "name",
                                        "",
                                },
                                {
                                        "summoner_level",
                                        12312,
                                },
                        },
                        " AND name = $1 AND summoner_level = $2",
                },
        } {
                query, _ := build("", 0, test.filters...)
                assert.Equal(t, test.want, query)
        }

}
