package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/roshdyosf/notificationSys/database"
	"github.com/roshdyosf/notificationSys/pkg/provider"
	"github.com/roshdyosf/notificationSys/repository"
	"github.com/roshdyosf/notificationSys/worker"
)

type App struct {
	router http.Handler
	db *sql.DB
	worker *worker.Worker
}



func New()*App{
	db,err := database.ConnectDB()
	if err!= nil{
		log.Fatalf("Database connection failed: %v", err )
	}


	repo := repository.NewPostgresNotificationRepo(db)
	emailProv := provider.NewMockEmailProvider()


	notifWorker := worker.NewWorker(repo, emailProv)

	app:= &App{

		db:db,
		worker: notifWorker,
		
	}
	app.router = app.loadRoutes(repo)
	return app

}

func (a *App) Start(ctx context.Context)error{
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" 
	}
	server:= &http.Server{
		Addr:  ":" + port,
		Handler: a.router,
	}
	go a.worker.Start(ctx)

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		return fmt.Errorf("server error: %w", err)

	case <-ctx.Done():
		log.Println("Shutting down gracefully...")

		defer a.db.Close()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("forced shutdown: %w", err)
		}

		log.Println("Server stopped cleanly.")
		return nil
	}
}