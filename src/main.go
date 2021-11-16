package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const listenerPort = ":8080"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	server := &http.Server{
		Addr: listenerPort,
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
