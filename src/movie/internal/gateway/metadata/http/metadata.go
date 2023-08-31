package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"microservices/metadata/pkg/model"
	"microservices/movie/internal/gateway"
	"microservices/pkg/discovery"
	"net/http"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addresses, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + addresses[rand.Intn(len(addresses))] + "/metadata"
	log.Printf("Calling metadata service. Request: GET " + url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	values := request.URL.Query()
	values.Add("id", id)
	request.URL.RawQuery = values.Encode()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", response)
	}
	var v *model.Metadata
	if err := json.NewDecoder(response.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, err
}
