package repository

import (
	"context"
	"database/sql"

	"github.com/roshdyosf/notificationSys/model"
)


type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
}

type postgresNotificationRepo struct {
	db *sql.DB
}

func NewPostgresNotificationRepo(db *sql.DB) NotificationRepository {
	return &postgresNotificationRepo{db: db}
}


func (r *postgresNotificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, message, status, is_read, retry_count, created_at, updated_at)
		VALUES ($1, $2, $3, 'Pending', false, 0, NOW(), NOW())
		RETURNING id, status, is_read, retry_count, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx, 
		query, 
		notif.UserID, 
		notif.Type, 
		notif.Message,
	).Scan(
		&notif.ID, 
		&notif.Status, 
		&notif.IsRead, 
		&notif.RetryCount, 
		&notif.CreatedAt, 
		&notif.UpdatedAt,
	)
}