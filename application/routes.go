package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/roshdyosf/notificationSys/handler"
)

func (a *App)loadRoutes() *chi.Mux{
	
	router:= chi.NewRouter()

	router.Use(middleware.Logger)
	
	router.Get("/", func(w http.ResponseWriter , r *http.Request){

		w.WriteHeader(http.StatusOK)
	})

	router.Route("/notification",a.loadNotificationRoutes)
return router
}

func (a *App) loadNotificationRoutes(router chi.Router) {
	
notifHandler:=handler.NewNotificationHandler(a.db)
router.Post("/",notifHandler.Create)
router.Get("/",notifHandler.List)
router.Get("/{id}",notifHandler.MarkAsRead)
}

