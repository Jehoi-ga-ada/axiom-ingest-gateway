package usecase

import (
	"sync"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/infrastructure"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	shared "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/domain"
	v1 "github.com/Jehoi-ga-ada/axiom-schema/gen/go/v1"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var eventPool = sync.Pool{
    New: func() interface{} { return &v1.Event{} },
}

type eventIngester struct {
	logger *zap.Logger
	dispatcher infrastructure.EventDispatcher
}

func NewEventIngester(logger *zap.Logger, dispatcher infrastructure.EventDispatcher) EventIngester {
	return &eventIngester{
		logger: logger,
		dispatcher: dispatcher,
	}
}

func (u *eventIngester) Execute(ctx *fasthttp.RequestCtx, req dto.CreateEventRequest) (string, error) {
	tenantID, ok := ctx.UserValue("tenant_id").(shared.TenantID)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	tenantIDstr := string(tenantID)

	e, err := domain.NewEvent(tenantIDstr, req.Type, req.Timestamp, req.RawBody)
	if err != nil {
		u.logger.Error("failed to create event with error",
			zap.Error(err),
			zap.String("event_type", req.Type),
		)
		return "", err
	}

	err = e.IsValid()
	if err != nil {
		return "", err
	}

	pb := eventPool.Get().(*v1.Event)
	pb.Reset()
	defer eventPool.Put(pb)

	pb.TenantId = tenantIDstr
	pb.Id = e.ID[:]
	pb.EventType = req.Type
    pb.Timestamp = timestamppb.New(e.Timestamp)
    pb.RawBody = e.RawBody

	data, err := proto.Marshal(pb)
	if err != nil {
		u.logger.Error("failed to serialize event",
			zap.Error(err),
		)
		return "", err
	}

	if err := u.dispatcher.Enqueue(data); err != nil {
		u.logger.Error("dispatcher queue full", zap.Error(err))
		return "", err
	}

	return e.ID.String(), nil
}