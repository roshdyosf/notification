package handler

import (
	"database/sql"
	"fmt"
	"net/http"
)

type NotificationHandler struct{
	DB *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

func (n *NotificationHandler) Create(w http.ResponseWriter,r *http.Request){
fmt.Println("Create notification")

}
func (n *NotificationHandler) List(w http.ResponseWriter,r *http.Request){
	fmt.Println("List notification")
}
func (n *NotificationHandler) MarkAsRead(w http.ResponseWriter,r *http.Request){
	fmt.Println("touch notification")
}

