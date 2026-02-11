package domain

import "errors"

var (
	ErrInvalidType         = errors.New("invalid event type")
	ErrInvalidTimestamp	   = errors.New("invalid timestamp")
	ErrInvalidPayload  	   = errors.New("invalid payload")
	ErrEventTooLarge       = errors.New("event too large")
	ErrDispatchQueueIsFull = errors.New("dispatcher queue saturated")
)
