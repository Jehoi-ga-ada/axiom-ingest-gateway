package http

import (
	"encoding/json"
	"net/http"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/usecase"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/dto"
	u "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	ei usecase.EventIngester	
}

func NewEventHandler(ei usecase.EventIngester) *EventHandler {
	return &EventHandler{
		ei: ei,
	}
}

func (h *EventHandler) Register(r chi.Router) {
	r.Group(func (r chi.Router) {
		r.Post("/", h.NewEvent)
	})
}

func (h *EventHandler) NewEvent(w http.ResponseWriter, r *http.Request) {
    r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	
	req := dto.CreateEventRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.BadRequest(w, err.Error())
		return
	}

	eventID, err := h.ei.Execute(r.Context(), req)
	if err != nil {
		u.BadRequest(w, err.Error())
		return
	}

	u.Created(w, dto.CreateEventResponse{
		EventID: eventID,
	})
}