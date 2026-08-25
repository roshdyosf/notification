package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
func (h *NotificationHandler) List(w http.ResponseWriter,r *http.Request){
	fmt.Println("List notification")
}
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter,r *http.Request){
	fmt.Println("touch notification")
}

