package dto

import (
	"encoding/json"
	"time"
)

type CreateEventRequest struct {
	Type string `json:"type" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
	RawBody json.RawMessage `json:"payload" binding:"required"`
}