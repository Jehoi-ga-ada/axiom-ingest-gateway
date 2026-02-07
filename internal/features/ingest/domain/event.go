package domain

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type EventType string
type Payload map[string]any

type Event struct {
	ID ulid.ULID
	Type EventType
	Timestamp time.Time
	RawBody []byte
}

func (e EventType) IsValid() bool {
	if len(e) == 0 || len(e) > 64 {
		return false
	}
	for _, r := range e {
		if !(r == '_' || r == '-' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func NewEvent(eventType string, timestamp time.Time, rawBody []byte) (*Event, error) {
    entropy := rand.Reader
    id, err := ulid.New(ulid.Timestamp(timestamp), entropy)
    if err != nil {
        return nil, err
    }

    return &Event{
        ID:        id,
        Type:      EventType(eventType),
        Timestamp: timestamp,
        RawBody:   rawBody,
    }, nil
}