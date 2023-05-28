package main

import (
	"log"
	"microservices/movie/internal/controller/movie"
	metadataGateway "microservices/movie/internal/gateway/metadata/http"
	ratingGateway "microservices/movie/internal/gateway/rating/http"
	httpHandler "microservices/movie/internal/handler/http"
	"net/http"
)

func main() {
	log.Printf("Starting the movie service")
	metadataGw := metadataGateway.New("localhost:8081")
	ratingGw := ratingGateway.New("localhost:8082")
	controller := movie.New(ratingGw, metadataGw)
	handler := httpHandler.New(controller)
	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
