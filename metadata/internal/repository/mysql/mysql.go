package mysql

import (
	"context"
	"database/sql"
	"errors"
	"movieapp/metadata/internal/repository"
	"movieapp/metadata/pkg/model"
	"os"
)

// Repository defines a MySQL-based movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New(db *sql.DB) (*Repository, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD") // TODO: implement better solution.
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@/"+dbName)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Get retrieves movie metadata by id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT title, description, director FROM metadata WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
	}
	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, m *model.Metadata) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)",
		id, m.Title, m.Description, m.Director)
	return err
}
