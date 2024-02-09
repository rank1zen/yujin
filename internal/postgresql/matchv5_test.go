package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectMatch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := NewDockerResource(t)

	pool, err := BackoffRetryPool(ctx, addr)
	require.NoError(t, err)

	err = Migrate(ctx, pool)
	require.NoError(t, err)

	db := newMatchV5Query(pool)

	for _, test := range []struct {
		arg  Match
		want string
	}{
		{
			arg: Match{Id: "JOEMAMA"},
			want: "JOEMAMA",
		},
	} {
		id, err := db.InsertMatch(ctx, &test.arg)
		assert.NoError(t, err)

		assert.Equal(t, test.want, id)
	}
}
