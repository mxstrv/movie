package memory

import (
	"context"
	"movieapp/metadata/internal/repository"
	"movieapp/metadata/pkg/model"
	"sync"
)

// Repository defines movie metadata in memory.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get method returns movie metadata by movie id
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata by id.
func (r *Repository) Put(_ context.Context, id string, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = m
	return nil
}
