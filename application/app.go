package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/roshdyosf/notificationSys/database"
)

type App struct {
	router http.Handler
	db *sql.DB
}



func New()*App{
	db,err := database.ConnectDB()
	if err!= nil{
		log.Fatalf("Database connection failed: %v", err )
	}

	app:= &App{

		db:db,
		
	}
	app.router= app.loadRoutes()
return app

}

func (a *App) Start(ctx context.Context)error{
port := os.Getenv("PORT")
	if port == "" {
		port = "4000" 
	}
Server:= &http.Server{
	Addr:  ":" + port,
	Handler: a.router,
}

defer a.db.Close()


err :=Server.ListenAndServe()

if err != nil {

	return fmt.Errorf("failed to start the server: %w",err)

}
return nil
}