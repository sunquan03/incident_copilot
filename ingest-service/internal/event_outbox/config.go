package event_outbox

import (
	"time"

	"github.com/sunquan03/ingest-service/internal/models"
)

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
	TopicFunc func(msg models.EventOutboxMessage) string
	// default id
	KeyFunc func(msg models.EventOutboxMessage) string
}

func DefaultRelayConfig() RelayConfig {
	return RelayConfig{
		PollInterval:    500 * time.Millisecond,
		BatchSize:       100,
		MaxRetries:      5,
		BaseBackoff:     5 * time.Second,
		DeleteOnPublish: false,
		TopicFunc:       func(m models.EventOutboxMessage) string { return m.AggregateType + "." + m.EventType },
		KeyFunc:         func(m models.EventOutboxMessage) string { return m.AggregateID },
	}
}

type Metrics struct {
	Published    func(aggregateType, eventType string)
	Failed       func(aggregateType, eventType string, err error)
	Retried      func(aggregateType, eventType string)
	BatchFetched func(n int)
}
