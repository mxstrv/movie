package movie

import (
	"context"
	"errors"
	"movieapp/movie/internal/gateway"

	metadata "movieapp/metadata/pkg/model"
	moviedetails "movieapp/movie/pkg/model"
	"movieapp/rating/pkg/model"
)

var ErrNotFound = errors.New("movie not found")

type RatingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadata.Metadata, error)
}

// Controller defines a movie service controller.
type Controller struct {
	ratingGateway   RatingGateway
	metadataGateway metadataGateway
}

// New creates a new controller.
func New(ratingGateway RatingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway: ratingGateway, metadataGateway: metadataGateway}
}

// Get returns a movie details including the aggregated rating and movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*moviedetails.MovieDetails, error) {
	getMetadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &moviedetails.MovieDetails{Metadata: *getMetadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, model.RecordID(id), model.RecordTypeMovie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Just proceed in this case
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}
	return details, nil
}
