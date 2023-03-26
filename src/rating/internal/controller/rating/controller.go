package rating

import (
	"context"
	"errors"
	"microservices/rating/internal/repository"
	model "microservices/rating/pkg"
)

// ErrNotFound is returned when no ratings are found for a record.
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, id model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, id model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller
func New(repo ratingRepository) *Controller {
	return &Controller{repo: repo}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregatedRating(ctx context.Context, id model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, id, recordType)
	if err != nil {
		if err == repository.ErrNotFound {
			return 0, ErrNotFound
		} else {
			return 0, err
		}
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record.
func (c *Controller) PutRating(ctx context.Context, id model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, id, recordType, rating)
}
