package config

import (
	ucEvent "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/usecase"
	eventHttp "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/delivery/http"
	"github.com/fasthttp/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Router *router.Router
	Viper *viper.Viper
	Logger *zap.Logger
}

func NewApp(config *Config) {
	eventIngester := ucEvent.NewEventIngester(config.Logger)
	eventHandler := eventHttp.NewEventHandler(eventIngester)

	v1 := config.Router.Group("/api/v1")

	// --- Events ---
	events := v1.Group("/events")
	eventHandler.Register(events)

	
}