package database

import "time"

type MatchFrame struct {
        Frame time.Duration
}

// Some frame to store a players stats at a point in time.
type ParticipantFrame struct {
        Frame time.Duration

        Puuid string
        // some stats like damage and gold and cs
}

// Something struct to store an event: buying item, skill level up 
type EventFrame struct {
        Frame time.Duration

        Type any // something happened hear 
}
