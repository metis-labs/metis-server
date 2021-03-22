package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"oss.navercorp.com/metis/metis-server/server/database"
)

type Client struct {
	client *mongo.Client
}

const (
	uri         = "mongodb://localhost:27017"
	dbName      = "metis"
	dialTimeout = 10
)

func NewClient() *Client {
	return &Client{}
}

func (d *Client) Dial(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, dialTimeout*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("Could not connect to MongoDB: %s\n", err.Error())
		return err
	}
	fmt.Println("Connected to MongoDB")

	d.client = client
	return nil
}

func (d *Client) CreateModel(ctx context.Context, name string) (*database.Model, error) {
	result, err := d.client.Database(dbName).Collection("models").InsertOne(ctx, bson.M{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return &database.Model{
		ID:   database.ID(result.InsertedID.(primitive.ObjectID).Hex()),
		Name: name,
	}, nil
}

func (d *Client) CreateProject(ctx context.Context, name string) (*database.Project, error) {
	now := time.Now()
	result, err := d.client.Database(dbName).Collection("projects").InsertOne(ctx, bson.M{
		"name":       name,
		"created_at": now,
	})
	if err != nil {
		return nil, err
	}

	return &database.Project{
		ID:        database.ID(result.InsertedID.(primitive.ObjectID).Hex()),
		Name:      name,
		CreatedAt: now,
	}, nil
}

func (d *Client) ListProjects(ctx context.Context) ([]*database.Project, error) {
	cursor, err := d.client.Database(dbName).Collection("projects").Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Print(err)
		}
	}()

	var projects []*database.Project
	for cursor.Next(ctx) {
		var project database.Project
		idHolder := struct {
			ID primitive.ObjectID `bson:"_id"`
		}{}
		if err := cursor.Decode(&idHolder); err != nil {
			return nil, err
		}
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		project.ID = database.ID(idHolder.ID.Hex())
		projects = append(projects, &project)
	}

	return projects, nil
}

// Close closes the client, releasing any open resources.
func (d *Client) Close(ctx context.Context) error {
	if err := d.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
