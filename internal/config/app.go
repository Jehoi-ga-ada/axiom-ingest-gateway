package config

import (
	"time"

	eventInfra "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/application/infrastructure"
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

func NewApp(config *Config) eventInfra.EventDispatcher {
	dispatcherConfig := eventInfra.DispatcherConfig{
		BatchSize: 1000,
		FlushInterval: 10 * time.Millisecond,
		MaxWorkers: 10,
		MaxSenders: 6,
		QueueSize: 50000,
		BufferMaxSize: 1024 * 1024,
		TargetAddr: "127.0.0.1:8080",
		// TargetAddr: "/tmp/axiom.sock",
		WriteTimeout: 5 * time.Second,
	}
	eventDispatcher := eventInfra.NewTCPDispatcher(config.Logger, dispatcherConfig)
	// eventDispatcher := eventInfra.NewUDSDispatcher(config.Logger, dispatcherConfig)
	eventIngester := ucEvent.NewEventIngester(config.Logger, eventDispatcher)
	eventHandler := eventHttp.NewEventHandler(eventIngester)

	v1 := config.Router.Group("/api/v1")

	// --- Events ---
	events := v1.Group("/events")
	eventHandler.Register(events)

	return eventDispatcher
}