package routers

import (
	"github.com/fasthttp/router"
	"github.com/sunquan03/ingest-service/internal/handlers"
	"github.com/valyala/fasthttp"
)

type Router struct {
	router        *router.Router
	reqHandler    *handlers.Handler
	healthHandler *handlers.HealthHandler
}

func NewRouter(reqHandler *handlers.Handler, healthHandler *handlers.HealthHandler) *Router {
	return &Router{
		router:        router.New(),
		reqHandler:    reqHandler,
		healthHandler: healthHandler,
	}

}

func (r *Router) Setup() fasthttp.RequestHandler {

	apiV1 := r.router.Group("/api/v1")
	apiV1.POST("/alert", r.reqHandler.HandleAlert)
	apiV1.POST("/logdoc", r.reqHandler.HandleLogDoc)
	apiV1.POST("/incident", r.reqHandler.HandleIncident)
	apiV1.GET("/health", r.healthHandler.HandleHealth)

	return r.router.Handler
}
