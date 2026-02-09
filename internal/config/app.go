package config

import (
	ucEvent "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/usecase"
	eventHttp "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/http"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Router chi.Router
	Viper *viper.Viper
	Logger *zap.Logger
}

func NewApp(config *Config) {
	eventIngester := ucEvent.NewEventIngester(config.Logger)
	eventHandler := eventHttp.NewEventHandler(eventIngester)

	config.Router.Route("/api/v1", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			eventHandler.Register(r)
		})
	})
}