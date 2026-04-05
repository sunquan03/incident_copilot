package event_outbox

import (
	"context"
	"log"
	"time"

	"github.com/sunquan03/ingest-service/internal/brokers"
	"github.com/sunquan03/ingest-service/internal/models"
	"github.com/sunquan03/ingest-service/internal/repositories"
)

type EventOutboxRelay struct {
	repo     repositories.IEventOutboxRepository
	producer brokers.Producer
	cfg      RelayConfig
	metrics  Metrics
}

func NewEventOutboxRelay(
	repo repositories.IEventOutboxRepository,
	producer brokers.Producer,
	cfg RelayConfig,
) *EventOutboxRelay {
	if cfg.TopicFunc == nil {
		cfg.TopicFunc = DefaultRelayConfig().TopicFunc
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = DefaultRelayConfig().KeyFunc
	}
	return &EventOutboxRelay{repo: repo, producer: producer, cfg: cfg}
}

func (r *EventOutboxRelay) WithMetrics(m Metrics) *EventOutboxRelay {
	r.metrics = m
	return r
}

func (r *EventOutboxRelay) process(ctx context.Context, message models.EventOutboxMessage) error {
	topic := r.cfg.TopicFunc(message)
	key := r.cfg.KeyFunc(message)

	err := r.producer.SendMessage(topic, key, message.Payload)
	if err != nil {
		return err
	}

	if err = r.repo.MarkPublished(ctx, message.ID); err != nil {
		return err
	}

	if r.metrics.Published != nil {
		r.metrics.Published(message.AggregateType, message.EventType)
	}
	return nil
}

func (r *EventOutboxRelay) runCycle(ctx context.Context) error {
	messages, err := r.repo.FetchPendingBatch(ctx, r.cfg.BatchSize)
	if err != nil {
		return err // handle failure + 1
	}
	if len(messages) == 0 {
		return nil
	}

	if r.metrics.BatchFetched != nil {
		r.metrics.BatchFetched(len(messages))
	}
	for _, msg := range messages {
		if err := r.process(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (r *EventOutboxRelay) Run(ctx context.Context) error {
	log.Println("[event_outbox] starting")

	ticker := time.NewTicker(r.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[event_outbox] stopping")

			return ctx.Err()
		case <-ticker.C:
			if err := r.runCycle(ctx); err != nil {
				return err
			}

		}
	}

}
