package main

import (
	"context"
	h "guber/services/trip-service/internal/infrastructure/http"
	"guber/services/trip-service/internal/infrastructure/repository"
	"guber/services/trip-service/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("🚀 Starting Trip")

	inMemRepo := repository.NewInMemRepository()
	svc := service.NewService(inMemRepo)
	mux := http.NewServeMux()

	httpHandler := h.HttpHandler{Service: svc}
	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Starting listening on :8083")
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
			log.Printf("Server is shutting down due %v", sig)
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
