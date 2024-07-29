package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProfileSummary(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	name, err := ParseRiotName("orrange-na1")
	require.NoError(t, err)

	err = db.UpdateProfile(ctx, name)
	require.NoError(t, err)

	_, err = db.GetProfileSummary(ctx, name)
	assert.NoError(t, err)
}

func TestGetProfileMatchlist(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	name, err := ParseRiotName("orrange-na1")
	require.NoError(t, err)

	err = db.UpdateProfile(ctx, name)
	require.NoError(t, err)

	_, err = db.GetProfileMatchList(ctx, name, 0)
	assert.NoError(t, err)
}
