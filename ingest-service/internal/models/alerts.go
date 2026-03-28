package models

type Alert struct {
	SourceID   string         `json:"source_id"`
	SourceName string         `json:"source_name"`
	Message    string         `json:"message"`
	Labels     map[string]any `json:"labels"`
	CreatedAt  int64          `json:"created_at"`
}
