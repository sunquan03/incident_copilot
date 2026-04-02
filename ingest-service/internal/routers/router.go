package routers

import (
	"github.com/fasthttp/router"
	"github.com/sunquan03/ingest-service/internal/handlers"
)

type Router struct {
	router        *router.Router
	reqHandler    *handlers.Handler
	healthHandler *handlers.HealthHandler
}

func NewRouter(router *router.Router, reqHandler *handlers.Handler, healthHandler *handlers.HealthHandler) *Router {
	return &Router{
		router:        router,
		reqHandler:    reqHandler,
		healthHandler: healthHandler,
	}

}
