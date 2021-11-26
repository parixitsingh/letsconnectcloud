package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"server/filemanager"
)

const listenerPort = ":8080"

func main() {
	server := &http.Server{
		Addr:    listenerPort,
		Handler: filemanager.Routes(),
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	// Set up signal handling to exit cleanly if needed
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Errorf("shutting down server: %v", err)
		}
		<-errChan
	case <-errChan:
		os.Exit(1)
	}
}
