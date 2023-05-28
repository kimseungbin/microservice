package main

import (
	"log"
	"microservices/metadata/internal/controller/metadata"
	httphandler "microservices/metadata/internal/handler/http"
	"microservices/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the movie metadata servicew")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
