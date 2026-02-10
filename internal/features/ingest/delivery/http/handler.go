package http

import (
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/usecase"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	u "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/utils"
	"github.com/bytedance/sonic"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type EventHandler struct {
	ei usecase.EventIngester	
}

func NewEventHandler(ei usecase.EventIngester) *EventHandler {
	return &EventHandler{
		ei: ei,
	}
}

func (h *EventHandler) Register(r *router.Group) {
	r.POST("/", h.NewEvent)
}

func (h *EventHandler) NewEvent(ctx *fasthttp.RequestCtx) {
	if len(ctx.PostBody()) > 1024 * 1024 {
		u.RequestEntityTooLarge(ctx, "Entity is too large, only 1MB is allowed")
	}
	
	req := dto.CreateEventRequest{}

	if err := sonic.Unmarshal(ctx.PostBody(), &req); err != nil {
		u.BadRequest(ctx, err.Error())
		return
	}

	eventID, err := h.ei.Execute(ctx, req)
	if err != nil {
		u.BadRequest(ctx, err.Error())
		return
	}

	resp := dto.CreateEventResponse{EventID: eventID}

	u.Created(ctx, resp)
}