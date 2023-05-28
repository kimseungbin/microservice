package model

import (
	"microservices/metadata/pkg/model"
)

// MovieDetails includes movie metadata and its aggregated rating.
type MovieDetails struct {
	Rating   *float64 `json:"rating,omitempty"`
	Metadata model.Metadata
}
