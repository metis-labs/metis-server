package database

import (
	"context"
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidID = errors.New("invalid id")
	ErrNotFound  = errors.New("resource not found")
)

// ID represents ID of entity.
type ID string

// String returns a string representation of this ID.
func (id ID) String() string {
	return string(id)
}

// Bytes returns bytes of decoded hexadecimal string representation of this ID.
func (id ID) Bytes() []byte {
	decoded, err := hex.DecodeString(id.String())
	if err != nil {
		return nil
	}
	return decoded
}

// IDFromBytes returns ID represented by the encoded hexadecimal string from bytes.
func IDFromBytes(bytes []byte) ID {
	return ID(hex.EncodeToString(bytes))
}

// Database represents database which reads or saves Metis data.
type Database interface {
	Dial(ctx context.Context) error
	Close(ctx context.Context) error

	CreateModel(ctx context.Context, name string) (*Model, error)
	CreateProject(ctx context.Context, name string) (*Project, error)
	ListProjects(ctx context.Context) ([]*Project, error)
	UpdateProject(ctx context.Context, id string, name string) error
	DeleteProject(ctx context.Context, id string) error
}
