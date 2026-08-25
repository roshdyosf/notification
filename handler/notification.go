package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/roshdyosf/notificationSys/model"
	"github.com/roshdyosf/notificationSys/repository"
)

type NotificationHandler struct {
	repo repository.NotificationRepository
}

func NewNotificationHandler(repo repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

type CreateNotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

var allowedNotificationTypes = map[string]bool{
	"EMAIL": true,
	"SMS":   true,
	"PUSH":  true,
}

func isValidNotificationType(nType string) bool {
	return allowedNotificationTypes[nType]
}

func (h *NotificationHandler) Create(w http.ResponseWriter,r *http.Request){
	var req CreateNotificationRequest


	//validation first
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	if req.UserID == "" || req.Type == "" || req.Message == "" {
			http.Error(w, "user_id, type, and message are required", http.StatusUnprocessableEntity)
			return
		}
	if !isValidNotificationType(req.Type) {
			http.Error(w, "Invalid notification type. Allowed types: EMAIL, SMS, PUSH", http.StatusBadRequest)
			return
		}
	notif := &model.Notification{
			UserID:  req.UserID,
			Type:    req.Type,
			Message: req.Message,
		}

	if err := h.repo.Create(r.Context(), notif); err != nil {
			http.Error(w, "Failed to save notification: "+err.Error(), http.StatusInternalServerError)
			return
		}

	w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Notification created successfully",
			"data":    notif,
		})

}


func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	
userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	unreadOnly := r.URL.Query().Get("unread_only") == "true"

	notifications, total, err := h.repo.ListByUserID(r.Context(),userID,limit,offset,unreadOnly)

	if err != nil {
		http.Error(w, "Failed to fetch notifications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": notifications,
			"meta": map[string]interface{}{
				"total_items": total,
				"page":        page,
				"limit":       limit,
				"total_pages": (total + limit - 1) / limit,
			},
		})
}

func (h *NotificationHandler) Read(w http.ResponseWriter,r *http.Request){
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	notif, err := h.repo.MarkAsRead(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to mark notification as read: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Notification retrieved and marked as read successfully",
		"data":    notif,
	})
}
