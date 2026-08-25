package model

import "time"
type NotificationStatus string

const (
	StatusPending    NotificationStatus = "Pending"
	StatusProcessing NotificationStatus = "Processing"
	StatusSuccess    NotificationStatus = "Success"
	StatusFail       NotificationStatus = "Fail"
)


type Notification struct {
	ID          string             `json:"id"`
	UserID      string             `json:"user_id"`
	Type        string             `json:"type"`
	Message     string             `json:"message"`
	Status      NotificationStatus `json:"status"`
	IsRead      bool               `json:"is_read"`
	RetryCount  int                `json:"retry_count"`
	ErrorReason string             `json:"error_reason,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}