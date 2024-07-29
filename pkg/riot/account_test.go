package riot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountByRiotId(t *testing.T) {
	client := setup(t)

	ctx := testingContext(t)

	_, err := client.GetAccountByRiotId(ctx, "orrange", "na1")
	assert.NoError(t, err)
}
