package usecase

import (
	"context"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
)

type eventIngester struct {
	
}

func NewEventIngester() EventIngester {
	return &eventIngester{}
}

func (u *eventIngester) Execute(ctx context.Context, req dto.CreateEventRequest) (string, error) {
	e, err := domain.NewEvent(req.Type, req.Timestamp, req.RawBody)

	if err != nil {
		return "", err
	}

	return e.ID.String(), nil
}