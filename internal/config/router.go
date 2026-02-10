package config

import "github.com/fasthttp/router"

func NewRouter() *router.Router {
	r := router.New()

	return r
}