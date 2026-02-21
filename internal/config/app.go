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
		BatchSize: config.Viper.GetInt("BATCH_SIZE"),
		FlushInterval: time.Duration(config.Viper.GetInt("FLUSH_INTERVAL")) * time.Millisecond,
		MaxWorkers: config.Viper.GetInt("MAX_WORKERS"),
		MaxSenders: config.Viper.GetInt("MAX_SENDERS"),
		QueueSize: config.Viper.GetInt("QUEUE_SIZE"),
		BufferMaxSize: config.Viper.GetInt("BUFFER_MAX_SIZE") * 1024,
		TargetAddr: config.Viper.GetString("DISPATCHER_ADDR"),
		WriteTimeout: time.Duration(config.Viper.GetInt("WRITE_TIMEOUT")) * time.Second,
	}
	eventDispatcher := eventInfra.NewTCPDispatcher(config.Logger, dispatcherConfig)
	eventIngester := ucEvent.NewEventIngester(config.Logger, eventDispatcher)
	eventHandler := eventHttp.NewEventHandler(eventIngester)

	v1 := config.Router.Group("/api/v1")

	// --- Events ---
	events := v1.Group("/events")
	eventHandler.Register(events)

	return eventDispatcher
}