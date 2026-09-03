package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/roshdyosf/notificationSys/handler"
	"github.com/roshdyosf/notificationSys/repository"
)

func (a *App)loadRoutes(repo repository.NotificationRepository ) *chi.Mux{
	
	router:= chi.NewRouter()

	router.Use(middleware.Logger)
	
	//health check
	router.Get("/", func(w http.ResponseWriter , r *http.Request){
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api/v1/notifications",func(r chi.Router) {
		a.loadNotificationRoutes(r, repo)
	})
return router
}

func (a *App) loadNotificationRoutes(router chi.Router,repo repository.NotificationRepository) {

	notifHandler := handler.NewNotificationHandler(repo)
	router.Post("/", notifHandler.Create)
	router.Get("/", notifHandler.List)
	router.Get("/{id}", notifHandler.Read)
}