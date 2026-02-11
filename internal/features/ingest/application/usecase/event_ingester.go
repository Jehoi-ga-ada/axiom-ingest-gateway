package usecase

import (
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	v1 "github.com/Jehoi-ga-ada/axiom-schema/gen/go/v1"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type eventIngester struct {
	logger *zap.Logger
	dispatcher domain.EventDispatcher
}

func NewEventIngester(logger *zap.Logger) EventIngester {
	return &eventIngester{
		logger: logger,
	}
}

func (u *eventIngester) Execute(ctx *fasthttp.RequestCtx, req dto.CreateEventRequest) (string, error) {
	e, err := domain.NewEvent(req.Type, req.Timestamp, req.RawBody)

	if err != nil {
		u.logger.Error("failed to create event with error",
			zap.Error(err),
			zap.String("event_type", req.Type),
		)
		return "", err
	}

	pb := &v1.Event{
		Id: e.ID[:],
		EventType: string(e.Type),
		Timestsamp: timestamppb.New(e.Timestamp),
		RawBody: e.RawBody,
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		u.logger.Error("failed to serialize event",
			zap.Error(err),
		)
		return "", err
	}

	if err := u.dispatcher.Enqueue(ctx, data); err != nil {
		u.logger.Error("dispatcher queue full", zap.Error(err))
		return "", err
	}

	return e.ID.String(), nil
}