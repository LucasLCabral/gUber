package main

import (
	"context"
	"guber/shared/env"
	"guber/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var gRPCAddr = ":9092"

func main() {
	log.Println("🚀 Starting Driver Service!")
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	service := NewServive()

	// Connecting to RabbitMQ
	conn, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	// Starting the gRPC server
	gRPCServer := grpcserver.NewServer()
	NewGRPCHandler(gRPCServer, service)

	log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	gRPCServer.GracefulStop()
}
