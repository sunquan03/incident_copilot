package models

import (
	"fmt"
	"time"
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
