package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/roshdyosf/notificationSys/model"
)


type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	ListByUserID(ctx context.Context, userID string, limit, offset int, unreadOnly bool)([]model.Notification, int, error)
	MarkAsRead(ctx context.Context, id string) (*model.Notification, error)
	FetchPending(ctx context.Context, limit int) ([]model.Notification, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateRetry(ctx context.Context, id string, retryCount int, status string) error
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


func (r *postgresNotificationRepo) ListByUserID(ctx context.Context, userID string, limit, offset int, unreadOnly bool) ([]model.Notification, int, error) {
	baseQuery := `FROM notifications WHERE user_id = $1`
	args := []interface{}{userID}
	paramIdx := 2

	if unreadOnly {
		baseQuery += " AND is_read = false"
	}

	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, user_id, type, message, status, is_read, retry_count, created_at, updated_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		baseQuery, paramIdx, paramIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	notifications := []model.Notification{}


	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Message,
			&n.Status,
			&n.IsRead,
			&n.RetryCount,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
		}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil

}


func (r *postgresNotificationRepo) MarkAsRead(ctx context.Context,id string)(*model.Notification, error ){
	query := `
			UPDATE notifications
			SET is_read = true, updated_at = NOW()
			WHERE id = $1
			RETURNING id, user_id, type, message, status, is_read, retry_count, created_at, updated_at
		`
	var n model.Notification
	err:=r.db.QueryRowContext(ctx,query,id).Scan(
		&n.ID,
		&n.UserID,
		&n.Type,
		&n.Message,
		&n.Status,
		&n.IsRead,
		&n.RetryCount,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows 
			}
			return nil, err
		}
	return &n, nil
}
