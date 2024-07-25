package riot

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	AndrewPUUID       = "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	AndrewAsheMatchID = "NA1_5003179356"
)

func TestMatch_List(t *testing.T) {
	t.Parallel()

	ctx := testingContext(t)

	client := setup(t)

	ids, err := client.GetMatchHistory(ctx, AndrewPUUID, 0, 5)
	if assert.NoError(t, err) {
		assert.Len(t, ids, 5)
	}
}

func TestMatch_GetMatch(t *testing.T) {
	t.Parallel()

	ctx := testingContext(t)

	client := setup(t)

	want := &MatchInfo{
		GameCreation:       1717303471865,
		GameStartTimestamp: 1717303529206,
		GameEndTimestamp:   1717304694311,
		GameDuration:       1165,
		GameVersion:        "14.11.589.9418",
		PlatformId:         "NA1",
	}

	m, err := client.GetMatch(ctx, "NA1_5011055088")
	if assert.NoError(t, err) {
		assert.NotNil(t, m.Info.Participants)
		assert.NotNil(t, m.Info.Teams)

		assert.Equal(t, want.GameCreation, m.Info.GameCreation)
		assert.Equal(t, want.GameStartTimestamp, m.Info.GameStartTimestamp)
		assert.Equal(t, want.GameEndTimestamp, m.Info.GameEndTimestamp)
		assert.Equal(t, want.GameDuration, m.Info.GameDuration)
		assert.Equal(t, want.GameVersion, m.Info.GameVersion)
		assert.Equal(t, want.PlatformId, m.Info.PlatformId)
	}
}

func TestMatch_GetMatchTimeline(t *testing.T) {
	t.Parallel()

	ctx := testingContext(t)

	client := setup(t)

	m, err := client.GetMatchTimeline(ctx, AndrewAsheMatchID)
	if assert.NoError(t, err) {
		log.Print(len(m.Info.Frames))
	}
}
