package middleware

import (
	"runtime/debug"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/utils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func RecoveryMiddleware(log *zap.Logger) func (next fasthttp.RequestHandler) fasthttp.RequestHandler {
    return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
        return func(ctx *fasthttp.RequestCtx) {
            defer func() {
                if r := recover(); r != nil {
                    log.Error("recovered from panic",
                        zap.Any("error", r),
                        zap.ByteString("stack", debug.Stack()),
                        zap.ByteString("path", ctx.Path()),
                    )

                    utils.StatusInternalServerError(ctx, "internal server error")
                }
            }()
            next(ctx)
        }
    }
}