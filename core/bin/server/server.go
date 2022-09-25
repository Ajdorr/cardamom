package main

import (
	cfg "cardamom/core/config"
	"cardamom/core/events"
	"cardamom/core/router"
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

	// Init server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.C.Host, cfg.C.Port),
		Handler: router.Engine,
	}

	// Start server wait
	go runServer(server)
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server was forced to shutdown: %s", err)
	}

	shutdown()
}
