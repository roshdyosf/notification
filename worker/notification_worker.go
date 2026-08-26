package worker

import (
	"context"
	"log"
	"time"

	"github.com/roshdyosf/notificationSys/pkg/provider"
	"github.com/roshdyosf/notificationSys/repository"
)

type Worker struct {
	repo          repository.NotificationRepository
	emailProvider provider.EmailProvider
	maxRetries    int
}

func NewWorker(repo repository.NotificationRepository, emailProv provider.EmailProvider) *Worker {
	return &Worker{
		repo:          repo,
		emailProvider: emailProv,
		maxRetries:    3, 
	}
}


func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	log.Println("🚀 Notification worker started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Stopping notification worker...")
			return
		case <-ticker.C:
			w.processNotifications(ctx)
		}
	}
}

func (w *Worker) processNotifications(ctx context.Context) {
	notifications, err := w.repo.FetchPending(ctx, 10)
	if err != nil  {
		log.Printf("[Worker Error] Failed to fetch pending notifications: %v", err)
		return
	}

	if len(notifications) == 0 {
				log.Printf("there is no pending notifications")
		return
	}

	for _, notif := range notifications {

		if err := w.repo.UpdateStatus(ctx, notif.ID, "Processing"); err != nil {
			log.Printf("[Worker Error] Failed to update status to Processing for ID %s: %v", notif.ID, err)
			continue
		}
		log.Printf("[Worker] Processing notification ID: %s | Type: %s", notif.ID, notif.Type)

		var sendErr error
		if notif.Type == "EMAIL" {
			sendErr = w.emailProvider.Send("mockemailfornow@mock.com",notif.Message)
		}

		if sendErr == nil {
			_ = w.repo.UpdateStatus(ctx, notif.ID, "Sent")
			log.Printf("[Worker]  Notification %s sent successfully!", notif.ID)
		} else {
			newRetryCount := notif.RetryCount + 1
			newStatus := "Pending" 

			if newRetryCount >= w.maxRetries {
				newStatus = "Failed" 
			}

			_ = w.repo.UpdateRetry(ctx, notif.ID, newRetryCount, newStatus)
			log.Printf("[Worker] Failed to send %s (Attempt %d/%d): %v", notif.ID, newRetryCount, w.maxRetries, sendErr)
		}
	}
}