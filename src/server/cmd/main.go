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

// listenerPort is the port on which store api are hosted
const listenerPort = ":8080"

func main() {
	// server is created
	server := &http.Server{
		Addr:    listenerPort,
		Handler: filemanager.Routes(),
	}

	errChan := make(chan error, 1)
	go func() {
		// listening on the server
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
		// In case of err from server exiting
		os.Exit(1)
	}
}
