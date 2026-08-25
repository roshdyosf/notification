package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/roshdyosf/notificationSys/handler"
)

func loadRoutes() *chi.Mux{
	
	router:= chi.NewRouter()

	router.Use(middleware.Logger)
	
	router.Get("/", func(w http.ResponseWriter , r *http.Request){

		w.WriteHeader(http.StatusOK)
	})

	router.Route("/notification",loadNotificationRoutes)
return router
}
func loadNotificationRoutes(router chi.Router){
	
notificationHandler := &handler.Notification{}
router.Post("/",notificationHandler.Create)
router.Get("/",notificationHandler.List)
router.Get("/{id}",notificationHandler.GetById)
}

