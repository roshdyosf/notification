package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/roshdyosf/notificationSys/application"
)

func main() {
	app := application.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx); err != nil {
		fmt.Println("Error running app:", err)
	}
}
