package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateSummoner(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	err := db.updateSummoner(ctx, RiotPuuid(AndrewPUUID))
	require.NoError(t, err)

}
