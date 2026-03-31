package repositories

import (
	"context"
	"time"

	"github.com/sunquan03/ingest-service/internal/database"
	"github.com/sunquan03/ingest-service/internal/models"
)

type IEventOutboxRepository interface {
	FetchPendingBatch(ctx context.Context, limit int) ([]models.EventOutboxMessage, error)
	MarkPublished(ctx context.Context, id int64) error
	ScheduleRetry(ctx context.Context, id int64, retryCount int, retryAt time.Time, err error) error
	MarkFailed(ctx context.Context, id int64, err error) error
}
type EventOutboxRepository struct {
	db *database.DB
}

func NewEventOutboxRepository(db *database.DB) *EventOutboxRepository {
	return &EventOutboxRepository{db: db}
}

func (repo *EventOutboxRepository) FetchPendingBatch(ctx context.Context, limit int) ([]models.EventOutboxMessage, error) {
	query := `select id, aggregate_type, aggregate_id, event_type, payload, retry_count, created_at
        from event_outbox
        where status = $1
          and (next_retry_at is NULL or next_retry_at <= NOW())
        order by created_at
        LIMIT $2
        FOR UPDATE SKIP LOCKED`

	rows, err := repo.db.QueryContext(ctx, query, models.StatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.EventOutboxMessage

	for rows.Next() {
		var msg models.EventOutboxMessage
		if err := rows.Scan(&msg.ID, &msg.AggregateType, &msg.AggregateID, &msg.EventType,
			&msg.Payload, &msg.RetryCount, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (repo *EventOutboxRepository) MarkPublished(ctx context.Context, id int64) error {
	query := `update event_outbox set 
                status = $1,
            	published_at = now(),
            	error_message = NULL,
            	next_retry_at = NULL
              where id = $2`

	_, err := repo.db.Exec(query, models.StatusPublished, id)
	return err
}

func (repo *EventOutboxRepository) ScheduleRetry(ctx context.Context, id int64, retryCount int, retryAt time.Time, errMsg string) error {
	query := `update event_outbox set 
                status = $1,
                retry_count = $2,
                next_retry_at = $3,
            	error_message = $4
              where id = $5`

	_, err := repo.db.ExecContext(ctx, query, models.StatusPending, retryCount, retryAt, errMsg, id)
	return err
}

func (repo *EventOutboxRepository) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	query := `update event_outbox set 
                status = $1,
            	error_message = $2,
            	next_retry_at = NULL
              where id = $3`

	_, err := repo.db.Exec(query, models.StatusFailed, errMsg, id)
	return err
}
