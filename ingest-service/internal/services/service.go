package services

import (
	"sunquan03/ingest-service/internal/brokers"
	"sunquan03/ingest-service/internal/repositories"
)

type Service struct {
	repo     *repositories.Repository
	producer brokers.Producer
}

func NewService(repo *repositories.Repository, producer brokers.Producer) Service {
	return Service{repo: repo, producer: producer}
}
