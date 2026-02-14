package domain

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/domain"
)

type EventType string

type Event struct {
	ID ulid.ULID
	TenantID domain.TenantID
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

func (e Event) IsValid() error {
	if !e.Type.IsValid() {
		return ErrInvalidType
	}

	// TODO: Add checking for other field

	return nil
}

func NewEvent(tenantID, eventType string, timestamp time.Time, rawBody []byte) (*Event, error) {
    entropy := rand.Reader
    id, err := ulid.New(ulid.Timestamp(timestamp), entropy)

    if err != nil {
        return nil, err
    }

    event := &Event{
		TenantID:  domain.TenantID(tenantID),
        ID:        id,
        Type:      EventType(eventType),
        Timestamp: timestamp,
        RawBody:   rawBody,
    }

	err = event.IsValid()

	if err != nil {
		return nil, err
	}

	return event, nil
}