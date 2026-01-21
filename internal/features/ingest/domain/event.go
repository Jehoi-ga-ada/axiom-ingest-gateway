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