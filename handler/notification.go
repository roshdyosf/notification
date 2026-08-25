package handler

import (
	"fmt"
	"net/http"
)

type Notification struct{}

func (n *Notification) Create(w http.ResponseWriter,r *http.Request){
fmt.Println("Create notification")

}
func (n *Notification) List(w http.ResponseWriter,r *http.Request){
	fmt.Println("List notification")
}
func (n *Notification) GetById(w http.ResponseWriter,r *http.Request){
	fmt.Println("touch notification")
}

