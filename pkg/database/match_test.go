package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	AndrewMatchID = "NA1_5011055088"
	AndrewPUUID   = "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
)

func TestEnsureMatchlist(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	err := db.ensureMatchlist(ctx, AndrewPUUID, 0, 1)
	assert.NoError(t, err)

	matches, err := db.getProfileMatchlist(ctx, AndrewPUUID, 0)
	assert.NoError(t, err)

	log.Print(matches[0])
}
