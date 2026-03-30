package event_outbox

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/sunquan03/ingest-service/internal/brokers"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPublished Status = "published"
	StatusFailed    Status = "failed"
)

type EventOutboxMessage struct {
	ID            int64
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	RetryCount    int
	CreatedAt     time.Time
}

func (m EventOutboxMessage) defaultTopic() string {
	return m.AggregateType + "." + m.EventType
}

func (m *EventOutboxMessage) defaultKey() string {
	return fmt.Sprintf("%d", m.ID)
}

type RelayConfig struct {
	// messages poll pending interval
	PollInterval time.Duration
	// max messages per single poll
	BatchSize int
	// max retry_count before failed status
	MaxRetries int
	// next_retry_at = now + BaseBackoff * 2^retry_count
	BaseBackoff time.Duration
	// delete published rows instead of setting status as published
	DeleteOnPublish bool
	// default "<aggregate_type>.<event_type>"
	TopicFunc func(msg EventOutboxMessage) string
	// default id
	KeyFunc func(msg EventOutboxMessage) string
}

func DefaultRelayConfig() RelayConfig {
	return RelayConfig{
		PollInterval:    500 * time.Millisecond,
		BatchSize:       100,
		MaxRetries:      5,
		BaseBackoff:     5 * time.Second,
		DeleteOnPublish: false,
		TopicFunc:       func(m EventOutboxMessage) string { return m.AggregateType + "." + m.EventType },
		KeyFunc:         func(m EventOutboxMessage) string { return m.AggregateID },
	}
}

type Metrics struct {
	Published    func(aggregateType, eventType string)
	Failed       func(aggregateType, eventType string, err error)
	Retried      func(aggregateType, eventType string)
	BatchFetched func(n int)
}

type EventOutboxRelay struct {
	db       *sql.DB
	producer brokers.Producer
	cfg      RelayConfig
	logger   *slog.Logger
	metrics  Metrics
}

func NewOutboxRelay(
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
