package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diagram struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Database interface {
	Dial(ctx context.Context) error
	Close(ctx context.Context) error

	CreateDiagram(ctx context.Context, name string) (*Diagram, error)
}
