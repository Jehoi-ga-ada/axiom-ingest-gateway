package utils

import (
	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/dto"
	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
)

func base(ctx *fasthttp.RequestCtx, code int, status string, data, errs any) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)

	res := dto.WebResponse{
		Status: status,
		Data: data,
		Errs: errs,
	}

	payload, _ := sonic.Marshal(res)
	ctx.SetBody(payload)
}

// --- Success Helpers ---
func Created(ctx *fasthttp.RequestCtx, data any) {
	base(ctx, fasthttp.StatusCreated, "CREATED", data, nil)
}

// --- Client Errors Helpers ---
func BadRequest(ctx *fasthttp.RequestCtx, message string) {
	base(ctx, fasthttp.StatusBadRequest, "BAD_REQUEST", nil, message)
}

func RequestEntityTooLarge(ctx *fasthttp.RequestCtx, message string) {
	base(ctx, fasthttp.StatusRequestEntityTooLarge, "STATUS_ENTITY_TOO_LARGE", nil, message)
}

// --- Server Helpers ---
func StatusServiceUnavailable(ctx *fasthttp.RequestCtx, message string) {
	base(ctx, fasthttp.StatusServiceUnavailable, "STATUS_SERVICE_UNAVAILABLE", nil, message)
}

func StatusInternalServerError(ctx *fasthttp.RequestCtx, message string) {
	base(ctx, fasthttp.StatusInternalServerError, "STATUS_INTERNAL_SERVER_ERROR", nil, message)
}