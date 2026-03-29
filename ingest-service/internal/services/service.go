package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/sunquan03/ingest-service/internal/brokers"
	"github.com/sunquan03/ingest-service/internal/models"
	"github.com/sunquan03/ingest-service/internal/repositories"
)

type Service struct {
	repo     *repositories.Repository
	producer brokers.Producer
}

func NewService(repo *repositories.Repository, producer brokers.Producer) Service {
	return Service{repo: repo, producer: producer}
}

func (s *Service) CreateIncident(ctx context.Context, inc *models.Incident) error {
	var id string
	id = "inc_" + uuid.New().String()
	err := s.repo.CreateIncident(ctx, id, inc)
	if err != nil {
		return err
	}

	err = s.repo.CreateEventOutbox(ctx, "incident", id, "incident.created", inc)
	return err
}

func (s *Service) CreateLogDoc(ctx context.Context, logdoc *models.LogDoc) error {
	var id string
	id = "logdoc_" + uuid.New().String()
	err := s.repo.CreateLogDoc(ctx, id, logdoc)
	if err != nil {
		return err
	}

	err = s.repo.CreateEventOutbox(ctx, "log_doc", id, "log_doc.created", logdoc)
	return err
}

func (s *Service) CreateAlert(ctx context.Context, alert *models.Alert) error {
	var id string
	id = "alert_" + uuid.New().String()
	err := s.repo.CreateAlert(ctx, id, alert)
	if err != nil {
		return err
	}

	err = s.repo.CreateEventOutbox(ctx, "alert", id, "alert.received", alert)

	return err
}
