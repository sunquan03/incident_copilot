package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/sunquan03/ingest-service/internal/models"
	"github.com/sunquan03/ingest-service/internal/services"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /api/v1/alert
func (h *Handler) HandleAlert(reqCtx *fasthttp.RequestCtx) {
	reqBody := reqCtx.PostBody()

	var req models.Alert

	if err := json.Unmarshal(reqBody, &req); err != nil {
		h.sendError(reqCtx, fasthttp.StatusBadRequest, err)
	}

	ctx := context.Background()

	err := h.service.CreateAlert(ctx, &req)
	if err != nil {
		h.sendError(reqCtx, fasthttp.StatusInternalServerError, err)
	}

	reqCtx.Response.Header.SetContentType("application/json")
	reqCtx.Response.SetStatusCode(fasthttp.StatusOK)

}

// POST /api/v1/logdoc
func (h *Handler) HandleLogDoc(reqCtx *fasthttp.RequestCtx) {
	reqBody := reqCtx.PostBody()

	var req models.LogDoc

	if err := json.Unmarshal(reqBody, &req); err != nil {
		h.sendError(reqCtx, fasthttp.StatusBadRequest, err)
	}

	ctx := context.Background()

	err := h.service.CreateLogDoc(ctx, &req)
	if err != nil {
		h.sendError(reqCtx, fasthttp.StatusInternalServerError, err)
	}

	reqCtx.Response.Header.SetContentType("application/json")
	reqCtx.Response.SetStatusCode(fasthttp.StatusOK)

}

// POST /api/v1/incident
func (h *Handler) HandleIncident(reqCtx *fasthttp.RequestCtx) {
	reqBody := reqCtx.PostBody()

	var req models.Incident

	if err := json.Unmarshal(reqBody, &req); err != nil {
		h.sendError(reqCtx, fasthttp.StatusBadRequest, err)
	}

	ctx := context.Background()

	err := h.service.CreateIncident(ctx, &req)
	if err != nil {
		h.sendError(reqCtx, fasthttp.StatusInternalServerError, err)
	}

	reqCtx.Response.Header.SetContentType("application/json")
	reqCtx.Response.SetStatusCode(fasthttp.StatusOK)

}

// error
func (h *Handler) sendError(reqCtx *fasthttp.RequestCtx, statusCode int, err error) {
	resBody := map[string]string{"status": "error", "message": err.Error()}

	reqCtx.Response.Header.SetContentType("application/json; charset=utf-8")
	reqCtx.Response.SetStatusCode(statusCode)

	if err = json.NewEncoder(reqCtx).Encode(resBody); err != nil {
		log.Printf("failed to encode error body to json: %v", err)
	}
}
