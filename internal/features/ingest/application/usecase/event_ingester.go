package usecase

import (
	"context"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	"go.uber.org/zap"
)

type eventIngester struct {
	logger *zap.Logger
}

func NewEventIngester(logger *zap.Logger) EventIngester {
	return &eventIngester{
		logger: logger,
	}
}

func (u *eventIngester) Execute(ctx context.Context, req dto.CreateEventRequest) (string, error) {
	e, err := domain.NewEvent(req.Type, req.Timestamp, req.RawBody)

	if err != nil {
		u.logger.Fatal("Failed to create event with error",
			zap.String("error", err.Error()),
		)
		return "", err
	}

	return e.ID.String(), nil
}