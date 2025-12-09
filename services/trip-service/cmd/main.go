package main

import (
	h "guber/services/trip-service/internal/infrastructure/http"
	"guber/services/trip-service/internal/infrastructure/repository"
	"guber/services/trip-service/internal/service"
	"log"
	"net/http"
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
	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
