package middleware

import (
    "runtime/debug"
    "go.uber.org/zap"
    "github.com/valyala/fasthttp"
)

func RecoveryMiddleware(log *zap.Logger, next fasthttp.RequestHandler) fasthttp.RequestHandler {
    return func(ctx *fasthttp.RequestCtx) {
        defer func() {
            if r := recover(); r != nil {
                log.Error("recovered from panic",
                    zap.Any("error", r),
                    zap.ByteString("stack", debug.Stack()),
                    zap.ByteString("path", ctx.Path()),
                )

                ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
                ctx.SetContentType("application/json")
                ctx.SetBodyString(`{"status":"error","message":"internal server error"}`)
            }
        }()
        next(ctx)
    }
}