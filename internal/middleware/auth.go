package middleware

import (
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/auth/application/service"
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/utils"
	"github.com/valyala/fasthttp"
)

func AuthMiddleware(registry service.KeyRegistry) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	headerKey := []byte("X-API-Key")
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func (ctx *fasthttp.RequestCtx) {
			key := ctx.Request.Header.PeekBytes(headerKey)

			if len(key) == 0 {
				utils.StatusUnauthorized(ctx, "API key not found")
				return
			}

			meta, ok := registry.Get(key)
			if !ok || !meta.Active{
				utils.StatusUnauthorized(ctx, "API key not found")
				return
			}
			
			ctx.SetUserValue("auth_meta", meta)
        	next(ctx)
		}
	}
}