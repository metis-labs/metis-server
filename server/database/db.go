package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Database interface {
	Dial(ctx context.Context) error
	Close(ctx context.Context) error

	CreateModel(ctx context.Context, name string) (*Model, error)
}
