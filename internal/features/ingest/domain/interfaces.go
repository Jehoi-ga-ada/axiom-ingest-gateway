package domain

import "github.com/valyala/fasthttp"

type EventDispatcher interface {
	Enqueue(ctx *fasthttp.RequestCtx, data []byte) error
	Close()
}