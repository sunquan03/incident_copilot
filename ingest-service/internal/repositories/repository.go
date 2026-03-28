package repositories

import (
	"context"
	"sunquan03/ingest-service/internal/database"
	"sunquan03/ingest-service/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateAlert(ctx context.Context, id string, alert *models.Alert) error {
	query := `INSERT INTO alerts (id, source_id, source_name, message, created_at) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(query, id, alert.SourceID, alert.SourceName, alert.Message, alert.CreatedAt)
	return err
}

func (repo *Repository) CreateLogDoc(ctx context.Context, id string, logdoc *models.LogDoc) error {
	query := `INSERT INTO log_docs (id, source_id, source_name, title, content, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(query, id, logdoc.SourceID, logdoc.ServiceName, logdoc.Title, logdoc.Content, logdoc.CreatedAt)
	return err
}

func (repo *Repository) CreateIncident(ctx context.Context, id string, inc *models.Incident) error {
	query := `INSERT INTO incidents (id, source_id, source_name, message, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(query, id, inc.SourceID, inc.ServiceName, inc.Message, inc.CreatedAt)
	return err
}
