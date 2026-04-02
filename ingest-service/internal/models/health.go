package models

type HealthStatus string

const (
	StatusUp   HealthStatus = "up"
	StatusDown HealthStatus = "down"
)

type ComponentHealth struct {
	Status  HealthStatus `json:"status"`
	Latency string       `json:"latency,omitempty"`
	Error   string       `json:"error,omitempty"`
}

type HealthResponse struct {
	Status     HealthStatus      `json:"status"`
	Uptime     string            `json:"uptime"`
	Version    string            `json:"version"`
	Components []ComponentHealth `json:"components"`
}
