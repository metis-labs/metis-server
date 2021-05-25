package database

import (
	"context"
	"errors"

	"oss.navercorp.com/metis/metis-server/server/types"
)

var (
	// ErrInvalidID is returned when the given string is ID.
	ErrInvalidID = errors.New("invalid id")

	// ErrNotFound is returned when the requested resource cannot be found.
	ErrNotFound = errors.New("resource not found")
)

// Database represents database which reads or saves Metis data.
type Database interface {
	Dial(ctx context.Context) error
	Close(ctx context.Context) error

	CreateProject(ctx context.Context, name string) (*types.Project, error)
	FindProject(ctx context.Context, id types.ID) (*types.Project, error)
	ListProjects(ctx context.Context) ([]*types.Project, error)
	UpdateProject(ctx context.Context, id types.ID, name string) error
	DeleteProject(ctx context.Context, id types.ID) error
}
