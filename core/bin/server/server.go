package main

import (
	cfg "cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/source/events"
	"cardamom/core/source/router"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Unable to start server: %s\n", err)
	}
}

func shutdown() {
	events.Shutdown()
}

func getContext() context.Context {

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}

func main() {

	ctx := getContext()

	db.Connect()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.C.Server.Host, cfg.C.Server.Port),
		Handler: router.Engine,
	}

	go runServer(server)
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server was forced to shutdown: %s", err)
	}

	shutdown()
}
