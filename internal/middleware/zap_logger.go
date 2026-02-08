package config

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func ZapLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				requestID := middleware.GetReqID(r.Context())
			
				log.Info("request completed",
					zap.String("request_id", requestID), 
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.Duration("latency", time.Since(t1)),
					zap.Int("size", ww.BytesWritten()),
					zap.String("remote_ip", r.RemoteAddr),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}