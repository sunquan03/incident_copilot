package models

import "time"

type Incident struct {
	SourceID    string    `json:"source_id"`
	ServiceName string    `json:"service_name"`
	Message     string    `json:"message"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
