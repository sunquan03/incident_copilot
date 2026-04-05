package models

import "time"

type LogDoc struct {
	SourceID    string    `json:"source_id"`
	ServiceName string    `json:"service_name"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	SourceType  string    `json:"source_type"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
