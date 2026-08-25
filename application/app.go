package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}
func New()*App{
	app:= &App{
		router: loadRoutes(),
	}
return app

}


func (a *App) Start(ctx context.Context)error{

Server:= &http.Server{
	Addr: ":4000",
	Handler: a.router,
}

err :=Server.ListenAndServe()

if err != nil {

	return fmt.Errorf("failed to start the server: %w",err)

}
return nil
}