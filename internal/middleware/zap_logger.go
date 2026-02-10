package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func ZapLogger(log *zap.Logger) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			t1 := time.Now()

			next(ctx)

			log.Info("request completed",
				zap.Uint64("request_id", ctx.ID()),
				zap.ByteString("method", ctx.Method()),
				zap.ByteString("path", ctx.Path()),
				zap.Int("status", ctx.Response.StatusCode()),
				zap.Duration("latency", time.Since(t1)),
				zap.Int("size", len(ctx.Response.Body())),
				zap.String("remote_ip", ctx.RemoteAddr().String()),
			)
		}
	}
}