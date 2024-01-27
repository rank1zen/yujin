package postgresql

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUUID(t *testing.T) {
	_, err := newUUID("string")
	require.Error(t, err)

	_, err = newUUID("40e6215d-b5c6-4896-987c-f30f3678f608")
	require.NoError(t, err)
}

func TestNewTimestamp(t *testing.T) {
	r := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	pgtime := (r)

	ex := pgtype.Timestamp{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	assert.Equal(t, pgtime, ex)
}
