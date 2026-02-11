package usecase

import (
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	"github.com/valyala/fasthttp"
)

type EventIngester interface {
	Execute(ctx *fasthttp.RequestCtx, req dto.CreateEventRequest) (string, error)
}