package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func createServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Slow request started...")
		time.Sleep(8 * time.Second)
		log.Println("Slow request completed")
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func runServer(ctx context.Context, server *http.Server, timeout time.Duration) error {
	// Buffered channel to capture startup errors. Buffered so that go routing won't block when sending the error
	serverError := make(chan error, 1)

	go func() {
		log.Print("Starting server...")
		// 1. ListenAndServe is a blocking operation so it must be used in go routine
		// 2. ErrServerClosed is expected during graceful shutdown so we don't want to send it to the channel, only the
		// unexpected errors
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
		close(serverError)
	}()

	stop := make(chan os.Signal, 1)
	// Listen for incoming SIGTERM signal
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	// Server error
	case err := <-serverError:
		return err
	// SIGTERM Received
	case <-stop:
		log.Println("Shutdown signal received...")
	// Parent context ended
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	// Server shutdowns, stops accepting requests and gives timeout time to end existing ones
	if err := server.Shutdown(shutdownCtx); err != nil {
		// If there is an error with server shutdown we immediately close the listener and all active connections. It is
		// called only when timeout timesout preserving the graceful behaviour
		if closeErr := server.Close(); closeErr != nil {
			// If both shutdown and close fail, merge the errors
			return errors.Join(err, closeErr)
		}
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}

func getConfiguration() (string, time.Duration) {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	timeout := os.Getenv("TIMEOUT")
	timeoutInSeconds, err := strconv.Atoi(timeout)
	if err != nil {
		return port, 30 * time.Second
	}

	return port, time.Duration(timeoutInSeconds) * time.Second
}

func main() {
	port, timeout := getConfiguration()

	server := createServer(port)

	if err := runServer(context.Background(), server, timeout); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
