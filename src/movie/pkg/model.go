package model

import model "microservices/metadata/pkg"

// MovieDetails includes movie metadata and its aggregated rating.
type MovieDetails struct {
	Rating   *float64 `json:"rating,omitempty"`
	Metadata model.Metadata
}
