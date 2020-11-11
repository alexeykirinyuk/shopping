package main

import (
	"context"
	"fmt"
	"github.com/alexeykirinyuk/shopping/products/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	connectionString = "mongodb://localhost:27017"
)

func createMongoClient(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Print("Connected to mongodb")

	return client, nil
}

func main() {
	fmt.Println("Starting server on port :50051...")

	// 50051 is the default port for gRPC
	// Ideally we'd use 0.0.0.0 instead of localhost as well
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	mongoCtx := context.Background()
	db, err := createMongoClient(mongoCtx)
	defer func() {
		err := db.Disconnect(mongoCtx)
		if err != nil {
			log.Fatalf("Unable to disconnect database: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("Unable to connect database: %v", err)
	}

	server := grpc.NewServer()
	service := api.CreateProductService(db)
	api.RegisterProductServiceServer(server, service)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
