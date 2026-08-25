package repository

import (
	"context"
	"database/sql"

	"github.com/roshdyosf/notificationSys/model"
)


type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	ListByUserID(ctx context.Context,userID string)([]model.Notification,error)
	MarkAsRead(ctx context.Context, id string) (*model.Notification, error)
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


func(r *postgresNotificationRepo) ListByUserID(ctx context.Context , userID string)([]model.Notification, error){

query := `
		SELECT id, user_id, type, message, status, is_read, retry_count, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err !=nil{
		return nil,err
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
			return nil, err
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil

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
