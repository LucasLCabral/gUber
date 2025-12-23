package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"guber/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("🚀 Starting API Gateway")

	mux := http.NewServeMux()

	mux.HandleFunc("/trip/preview", enableCORS(handleTripPreview))
	mux.HandleFunc("/trip/start", enableCORS(handleTripStart))
	mux.HandleFunc("/ws/drivers", enableCORS(handleDriversWebSocket))
	mux.HandleFunc("/ws/riders", enableCORS(handleRidersWebSocket))


	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		{
			log.Printf("Error starting server: %v", err)
		}
	case sig := <-shutdown:
		{
			log.Printf("Server is shutting down due to %v", sig)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Could not complete graceful shutdown: %v", err)
				server.Close()
			}
			log.Println("Server stopped")
		}
	}
	
}
