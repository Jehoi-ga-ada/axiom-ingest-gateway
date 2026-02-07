package usecase

import (
	"context"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
)

type EventIngester interface {
	Execute(ctx context.Context, req dto.CreateEventRequest) (string, error)
}