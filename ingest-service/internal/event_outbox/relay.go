package event_outbox

import (
	"database/sql"
	"log/slog"

	"github.com/sunquan03/ingest-service/internal/brokers"
)

type EventOutboxRelay struct {
	db       *sql.DB
	producer brokers.Producer
	cfg      RelayConfig
	logger   *slog.Logger
	metrics  Metrics
}

func NewEventOutboxRelay(
	db *sql.DB,
	producer brokers.Producer,
	cfg RelayConfig,
	logger *slog.Logger,
) *EventOutboxRelay {
	if cfg.TopicFunc == nil {
		cfg.TopicFunc = DefaultRelayConfig().TopicFunc
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = DefaultRelayConfig().KeyFunc
	}
	return &EventOutboxRelay{db: db, producer: producer, cfg: cfg, logger: logger}
}

func (r *EventOutboxRelay) WithMetrics(m Metrics) *EventOutboxRelay {
	r.metrics = m
	return r
}
