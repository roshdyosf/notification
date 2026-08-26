package repository

import (
	"context"

	"github.com/roshdyosf/notificationSys/model"
)

func (r *postgresNotificationRepo) FetchPending(ctx context.Context, limit int) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, type, message, status, is_read, retry_count, created_at, updated_at
		FROM notifications
		WHERE status IN ('Pending', 'Failed') AND retry_count < 3
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.Status, &n.IsRead, &n.RetryCount, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *postgresNotificationRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE notifications SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *postgresNotificationRepo) UpdateRetry(ctx context.Context, id string, retryCount int, status string) error {
	query := `UPDATE notifications SET retry_count = $1, status = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, retryCount, status, id)
	return err
}