package config

import (
	"time"

	imid "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(imid.ZapLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return r
}