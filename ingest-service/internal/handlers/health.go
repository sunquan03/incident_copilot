package handlers

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/sunquan03/ingest-service/internal/database"
	"github.com/sunquan03/ingest-service/internal/models"
	"github.com/valyala/fasthttp"
)

type HealthHandler struct {
	db    *database.DB
	kafka sarama.Client
}

func NewHealthHandler(db *database.DB, kafka sarama.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		kafka: kafka,
	}
}

func (h *HealthHandler) checkKafka() models.ComponentHealth {
	start := time.Now()
	if h.kafka == nil {
		return models.ComponentHealth{Status: models.StatusDown, Error: "kafka not initialized"}
	}
	brokers := h.kafka.Brokers()
	if len(brokers) == 0 {
		return models.ComponentHealth{Status: models.StatusDown, Error: "no kafka brokers available"}
	}

	if err := h.kafka.RefreshMetadata(); err != nil {
		return models.ComponentHealth{Status: models.StatusDown, Error: err.Error()}
	}

	return models.ComponentHealth{
		Status:  models.StatusUp,
		Latency: time.Since(start).Round(time.Millisecond).String(),
	}
}

func (h *HealthHandler) checkPostgres() models.ComponentHealth {
	start := time.Now()
	if h.db == nil {
		return models.ComponentHealth{Status: models.StatusDown, Error: "postgres not initialized"}
	}
	ctx := context.Background()
	if err := h.db.Ping(ctx); err != nil {
		return models.ComponentHealth{Status: models.StatusDown, Error: err.Error()}
	}

	return models.ComponentHealth{
		Status:  models.StatusUp,
		Latency: time.Since(start).Round(time.Millisecond).String(),
	}
}

// api/v1/health
func (h *HealthHandler) HandleHealth(reqCtx *fasthttp.RequestCtx) {

}
