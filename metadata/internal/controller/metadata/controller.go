package metadata

import (
	"context"
	"errors"
	"movieapp/metadata/pkg/model"
)

// ErrNotFound is returned when requested record not found.
var ErrNotFound = errors.New("metadata not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// Controller defines a metadata controller.
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo: repo}
}

// Get returns metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	return res, nil
}
