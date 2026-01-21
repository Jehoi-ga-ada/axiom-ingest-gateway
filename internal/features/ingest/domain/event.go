package domain

import "time"

type EventType string
type Payload map[string]any

type Event struct {
	Type EventType
	Timestamp time.Time
	Payload Payload
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

func NewEvent(
	eventType string,
	timestamp time.Time,
	payload map[string]any,
	now time.Time,
) (*Event, error) {
	et := EventType(eventType)
	if !et.IsValid() {
		return nil, ErrInvalidType
	}

	if timestamp.IsZero() {
		return nil, ErrInvalidTimestamp
	}

	if timestamp.After(now.Add(5*time.Minute)) || timestamp.Before(now.Add(-24*time.Hour)) {
		return nil, ErrInvalidTimestamp
	}

	if payload == nil {
		return nil, ErrInvalidPayload
	}

	return &Event{
		Type: et,
		Timestamp: timestamp,
		Payload: payload,
	}, nil
}