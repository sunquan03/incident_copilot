package repositories

import (
	"context"

	"github.com/sunquan03/ingest-service/internal/database"
	"github.com/sunquan03/ingest-service/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateEventOutbox(ctx context.Context, aggrType, aggrId, eventType string, payload any) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	query := `INSERT INTO event_outbox (aggregate_type, aggregate_id, event_type, payload) 
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`
	_, err = tx.ExecContext(ctx, query, aggrType, aggrId, eventType, payload)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (repo *Repository) CreateAlert(ctx context.Context, id string, alert *models.Alert) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `INSERT INTO alerts (id, source_id, source_name, message, created_at) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, id, alert.SourceID, alert.SourceName, alert.Message, alert.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
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
