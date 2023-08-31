package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"microservices/movie/internal/controller/movie"
	metadataGateway "microservices/movie/internal/gateway/metadata/http"
	ratingGateway "microservices/movie/internal/gateway/rating/http"
	httpHandler "microservices/movie/internal/handler/http"
	"microservices/pkg/discovery"
	"microservices/pkg/discovery/consul"
	"net/http"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadataGateway.New(registry)
	ratingGateway := ratingGateway.New(registry)
	svc := movie.New(ratingGateway, metadataGateway)
	handler := httpHandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
